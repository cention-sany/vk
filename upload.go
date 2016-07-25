package vk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	VK_MAX_PHOTOS = 5
	VK_MAX_AUDIOS = 1
	VK_MAX_DOCS   = 1
	VK_MAX_VIDEOS = 1
	fileFieldName = "file"
)

type (
	UploadServer struct {
		UploadUrl string `json:"upload_url"`
		AlbumId   int    `json:"aid"`
		UserId    int    `json:"mid"`
		OwnerId   int    `json:"owner_id"`
		AccessKey string `json:"access_key"`
	}

	UploadAlbumPhotos struct {
		Server     int    `json:"server"`
		PhotosList string `json:"photos_list"`
		AID        int    `json:"aid"`
		Hash       string `json:"hash"`
	}

	UploadWallPhotos struct {
		Server int    `json:"server"`
		Photo  string `json:"photo"`
		Hash   string `json:"hash"`
	}

	UploadAudios struct {
		Redirect string `json:"redirect"`
		Server   int    `json:"server"`
		Audio    string `json:"audio"`
		Hash     string `json:"hash"`
	}

	UploadVideos struct {
		Size    int `json:"size"`
		VideoID int `json:"video_id"`
	}

	UploadDocs struct {
		File string `json:"file"`
	}
)

func (u *UploadAlbumPhotos) parseType() interface{} {
	return u
}

func (u *UploadWallPhotos) parseType() interface{} {
	return u
}

func (u *UploadAudios) parseType() interface{} {
	return u
}

func (u *UploadVideos) parseType() interface{} {
	return u
}

func (u *UploadDocs) parseType() interface{} {
	return u
}

func (u *UploadAlbumPhotos) useStream() bool {
	return true
}

func (u *UploadWallPhotos) useStream() bool {
	return true
}

func (u *UploadAudios) useStream() bool {
	return true
}

func (u *UploadVideos) useStream() bool {
	return false
}

func (u *UploadDocs) useStream() bool {
	return true
}

type uploader interface {
	uploadUrl(*Session) (string, error)
	field(int) string
	values() url.Values
	max() int
	parseType() interface{}
	useStream() bool
	postParse(*Session, []string) (json.RawMessage, error)
	format(r json.RawMessage) ([]string, error)
}

type baseUpload struct {
	mUp, f, u string
	v         url.Values
	limit, id int
}

func (u *baseUpload) uploadUrl(s *Session) (string, error) {
	var v UploadServer
	err := s.CallAPI(u.mUp, u.v, &v)
	if err != nil {
		return "", err
	}
	if v.UserId != 0 {
		u.id = v.UserId
	} else if v.OwnerId != 0 {
		u.id = v.OwnerId
	}
	//log.Println("debug: upload URL:", v.UploadUrl)
	return v.UploadUrl, nil
}

func (u *baseUpload) max() int {
	return u.limit
}

// return multipart field name
func (u *baseUpload) field(index int) string {
	if u.limit <= 1 {
		return u.f
	}
	return fmt.Sprint(u.f, strconv.Itoa(index))
}

func (u *baseUpload) values() url.Values {
	return u.v
}

// Photo (album)
const (
	photoAlbumUploadServer = "photos.getUploadServer"
	photoAlbumSave         = "photos.save"
)

type photoFormatter struct{}

func (photoFormatter) format(r json.RawMessage) ([]string, error) {
	var ps []Photo
	err := unmarshalToType(r, &ps)
	if err != nil {
		return nil, err
	}
	as := make([]string, len(ps))
	for i, p := range ps {
		as[i] = fmt.Sprint("photo", strconv.Itoa(p.OwnerId), "_", strconv.Itoa(p.Id))
	}
	return as, nil
}

type photoAlbumUpload struct {
	*baseUpload
	*UploadAlbumPhotos
	photoFormatter
}

func (u *photoAlbumUpload) postParse(s *Session, ns []string) (json.RawMessage, error) {
	var res json.RawMessage
	u.v.Set("photos_list", u.PhotosList)
	u.v.Set("server", strconv.Itoa(u.Server))
	u.v.Set("hash", u.Hash)
	if err := s.CallAPI(photoAlbumSave, u.v, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func getAlbumPhotoUploader(gid, a int) uploader {
	v := url.Values{}
	v.Set("album_id", strconv.Itoa(a))
	if gid != 0 {
		v.Set("group_id", strconv.Itoa(gid))
	}
	return &photoAlbumUpload{
		baseUpload: &baseUpload{
			mUp:   photoAlbumUploadServer,
			f:     fileFieldName,
			v:     v,
			limit: VK_MAX_PHOTOS,
		},
		UploadAlbumPhotos: &UploadAlbumPhotos{},
	}
}

// a is the album id
func newUserAlbumPhotoUploader(a int) uploader {
	return getAlbumPhotoUploader(0, a)
}

// gid is the group ID and a is the album id
func newCommunityAlbumPhotoUploader(gid, a int) uploader {
	return getAlbumPhotoUploader(gid, a)
}

// Photo (wall)
const (
	wallPhotoUploadServer = "photos.getWallUploadServer"
	wallPhotoSave         = "photos.saveWallPhoto"
)

type wallPhotoUpload struct {
	*baseUpload
	*UploadWallPhotos
	photoFormatter
}

func (u *wallPhotoUpload) postParse(s *Session, ns []string) (json.RawMessage, error) {
	var res json.RawMessage
	u.v.Set("photo", u.Photo)
	u.v.Set("server", strconv.Itoa(u.Server))
	u.v.Set("hash", u.Hash)
	if err := s.CallAPI(wallPhotoSave, u.v, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func getWallPhotoUploader(gid int) uploader {
	v := url.Values{}
	if gid != 0 {
		v.Set("group_id", strconv.Itoa(gid))
	}
	return &wallPhotoUpload{
		baseUpload: &baseUpload{
			mUp:   wallPhotoUploadServer,
			f:     fileFieldName,
			v:     v,
			limit: VK_MAX_PHOTOS,
		},
		UploadWallPhotos: &UploadWallPhotos{},
	}
}

func newUserWallPhotoUploader() uploader {
	return getWallPhotoUploader(0)
}

// gid is the group ID.
func newCommunityWallPhotoUploader(gid int) uploader {
	return getWallPhotoUploader(gid)
}

// Private message photos upload
const (
	pmPhotoUploadServer = "photos.getMessagesUploadServer"
	pmPhotoSave         = "photos.saveMessagesPhoto"
)

type pmPhotoUpload struct {
	*baseUpload
	*UploadWallPhotos
	photoFormatter
}

func (u *pmPhotoUpload) postParse(s *Session, ns []string) (json.RawMessage, error) {
	var res json.RawMessage
	u.v.Set("photo", u.Photo)
	u.v.Set("server", strconv.Itoa(u.Server))
	u.v.Set("hash", u.Hash)
	if err := s.CallAPI(pmPhotoSave, u.v, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// depend on current token user / group token
func getPMPhotoUploader(gid int) uploader {
	v := url.Values{}
	if gid != 0 {
		v.Set("group_id", strconv.Itoa(gid))
	}
	return &pmPhotoUpload{
		baseUpload: &baseUpload{
			mUp:   pmPhotoUploadServer,
			f:     fileFieldName,
			v:     v,
			limit: VK_MAX_PHOTOS,
		},
		UploadWallPhotos: &UploadWallPhotos{},
	}
}

// Audio
type audioUpload struct {
	*baseUpload
	*UploadAudios
}

func (u *audioUpload) postParse(s *Session, ns []string) (json.RawMessage, error) {
	var res json.RawMessage
	u.v.Set("audio", u.Audio)
	u.v.Set("server", strconv.Itoa(u.Server))
	u.v.Set("hash", u.Hash)
	if err := s.CallAPI("audio.save", u.v, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (u *audioUpload) format(r json.RawMessage) ([]string, error) {
	var a Audio
	err := unmarshalToType(r, &a)
	if err != nil {
		return nil, err
	}
	s := fmt.Sprint("audio", strconv.Itoa(a.Owner), "_", strconv.Itoa(a.Id))
	return []string{s}, nil
}

func getAudioUploader(gid int) uploader {
	v := url.Values{}
	if gid != 0 {
		v.Set("group_id", strconv.Itoa(gid))
	}
	return &audioUpload{
		baseUpload: &baseUpload{
			mUp:   "audio.getUploadServer",
			f:     fileFieldName,
			v:     v,
			limit: VK_MAX_AUDIOS,
		},
		UploadAudios: &UploadAudios{},
	}
}

func newUserAudioUploader() uploader {
	return getAudioUploader(0)
}

func newCommunityAudioUploader(gid int) uploader {
	return getAudioUploader(gid)
}

// Docs
const (
	wallDocUploadServer = "docs.getWallUploadServer"
	docUploadServer     = "docs.getUploadServer"
)

type docUpload struct {
	*baseUpload
	*UploadDocs
}

type docUnsaveType struct {
	Owner int `json:"owner"`
	Id    int `json:"id"`
}

func (u *docUpload) postParse(s *Session, ns []string) (json.RawMessage, error) {
	var res json.RawMessage
	u.v.Set("file", u.File)
	if len(ns) > 0 {
		u.v.Set("title", ns[0])
		u.v.Set("tags", ns[0])
	}
	err := s.CallAPI("docs.save", u.v, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (docUpload) format(r json.RawMessage) ([]string, error) {
	var docs []Doc
	err := unmarshalToType(r, &docs)
	if err != nil {
		return nil, err
	}
	as := make([]string, len(docs))
	for i, d := range docs {
		as[i] = fmt.Sprint("doc", strconv.Itoa(d.OwnerId), "_", strconv.Itoa(d.Id))
	}
	return as, nil
}

func getDocUploader(uurl string, gid int) uploader {
	v := url.Values{}
	if gid != 0 {
		v.Set("group_id", strconv.Itoa(gid))
	}
	return &docUpload{
		baseUpload: &baseUpload{
			mUp:   uurl,
			f:     fileFieldName,
			v:     v,
			limit: VK_MAX_DOCS,
		},
		UploadDocs: &UploadDocs{},
	}
}

func newUserWallDocUploader() uploader {
	return getDocUploader(wallDocUploadServer, 0)
}

func newCommunityWallDocUploader(gid int) uploader {
	return getDocUploader(wallDocUploadServer, gid)
}

func newUserDocUploader() uploader {
	return getDocUploader(docUploadServer, 0)
}

func newCommunityDocUploader(gid int) uploader {
	return getDocUploader(docUploadServer, gid)
}

// Videos
type videoUpload struct {
	*baseUpload
	*UploadVideos
	owner int
}

func (u *videoUpload) postParse(s *Session, ns []string) (json.RawMessage, error) {
	return nil, nil
}

func (u *videoUpload) format(r json.RawMessage) ([]string, error) {
	s := fmt.Sprint("video", fmt.Sprint(strconv.Itoa(u.id), "_", strconv.Itoa(u.VideoID)))
	return []string{s}, nil
}

func newVideoUploader(v url.Values, owner int) uploader {
	return &videoUpload{
		baseUpload: &baseUpload{
			mUp:   "video.save",
			f:     "video_file",
			v:     v,
			limit: VK_MAX_VIDEOS,
		},
		UploadVideos: &UploadVideos{},
		owner:        owner,
	}
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func writeParts(ps, ns []string, up uploader, w *multipart.Writer) error {
	var err error
	var i int
	var fieldName string
	n := 1
	psSize := len(ps)
	for i = 0; i < up.max() && i < psSize; i++ {
		if ps[i] == "" {
			continue
		}
		fieldName = up.field(n)
		err = func(name string) error {
			file, err := os.Open(ps[i])
			if err != nil {
				return err
			}
			defer file.Close()
			// h := make(textproto.MIMEHeader)
			// h.Set("Content-Disposition",
			// fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			// 	escapeQuotes(name), escapeQuotes(ns[i]))
			// h.Set("Content-Type", "video/mp4")
			// part, err := w.CreatePart(h)
			fileName := ns[i]
			if path.Ext(fileName) == "" {
				// avoid upload file because no dot extension
				fileName = fmt.Sprint(fileName, ".dat")
			}
			part, err := w.CreateFormFile(name, fileName)
			if err != nil {
				return err
			}
			if _, err = io.Copy(part, file); err != nil {
				return err
			}
			return nil
		}(fieldName)
		if err != nil {
			// report error
			break
		} else {
			n++
		}
	}
	if i < psSize {
		// this error can not be reported
	}
	return err
}

// sync name(s) base on path(s) and return the synced new names.
func syncNamesToPaths(ps, ns []string) []string {
	pSize := len(ps)
	nSize := len(ns)
	if pSize == nSize {
		return ns
	}
	if pSize < nSize {
		return ns[:pSize]
	}
	nns := make([]string, pSize)
	copy(nns, ns)
	for i := nSize; i < pSize; i++ {
		nns[i] = filepath.Base(ps[i])
	}
	return nns
}

func (s *Session) upload(ps, ns []string, up uploader) (json.RawMessage, error) {
	var (
		bb *bytes.Buffer
		pw *io.PipeWriter
		r  io.Reader
		w  io.Writer
	)
	uurl, err := up.uploadUrl(s)
	if err != nil {
		return nil, err
	}
	if up.useStream() {
		var pr *io.PipeReader
		pr, pw = io.Pipe()
		w = pw
		r = pr
	} else {
		bb = new(bytes.Buffer)
		w = bb
		r = bb
	}
	request, err := http.NewRequest("POST", uurl, r)
	if err != nil {
		return nil, err
	}
	writer := multipart.NewWriter(w)
	ns = syncNamesToPaths(ps, ns)
	if up.useStream() {
		go func() {
			err := writeParts(ps, ns, up, writer)
			if err == nil {
				if err = writer.Close(); err == nil {
					// no error
					if err = pw.Close(); err != nil {
						// this error can not be reported
					}
					return
				}
			} else if err = pw.CloseWithError(err); err != nil { // report multipart writing error
				// this error can not be reported
			}
		}()
	} else {
		if err = writeParts(ps, ns, up, writer); err != nil {
			return nil, err
		}
		if err = writer.Close(); err != nil {
			return nil, err
		}
		if bb != nil {
			request.ContentLength = int64(bb.Len())
		}
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	if err = d.Decode(up.parseType()); err != nil {
		return nil, err
	}
	return up.postParse(s, ns)
}

// UploadPhotosToAlbum upload photos to album - either community album or
// user album. ss are the path to the photo path. gid is the community id
// if community album is the destination and aid is the album id.
func (s *Session) UploadPhotosToAlbum(ss []string, gid, aid int) (json.RawMessage, error) {
	return s.upload(ss, nil, getAlbumPhotoUploader(gid, aid))
}

func (s *Session) UploadPhotosToWall(ss []string, gid int) (json.RawMessage, error) {
	return s.upload(ss, nil, getWallPhotoUploader(gid))
}

func (s *Session) UploadPhotosToPM(ss []string, gid int) (json.RawMessage, error) {
	return s.upload(ss, nil, getPMPhotoUploader(gid))
}

func (s *Session) UploadVideos(path string, v url.Values) error {
	_, err := s.upload([]string{path}, nil, newVideoUploader(v, 0))
	return err
}

// UploadDocsToWall only support upload to user wall at the moment. VK server
// disallow group docs upload. Only user token will work.
func (s *Session) UploadDocToWall(path string) (json.RawMessage, error) {
	return s.upload([]string{path}, nil, newUserWallDocUploader())
}

// UploadDocs only support upload to user docs at the moment. VK server
// disallow group docs upload. Only user token will work.
func (s *Session) UploadDoc(path string) (json.RawMessage, error) {
	return s.upload([]string{path}, nil, newUserDocUploader())
}

// UploadAudio only support upload to user audio at the moment. No API for
// group audio upload. gid is ignored.
func (s *Session) UploadAudio(path string, gid int) (json.RawMessage, error) {
	return s.upload([]string{path}, nil, getAudioUploader(gid))
}

func unmarshalToType(r json.RawMessage, v interface{}) error {
	b, err := r.MarshalJSON()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}

// ps are path(s) meanwhile ns are name(s)
func multiUploads(s *Session, ps, ns []string, u uploader) ([]string, error) {
	var r json.RawMessage
	var err error
	var paths, names, as []string
	ns = syncNamesToPaths(ps, ns)
	maxSize := u.max()
	for len(ps) > 0 {
		paths = ps
		names = ns
		if len(ps) > maxSize {
			paths = ps[:maxSize]
			names = ns[:maxSize]
			ps = ps[maxSize:]
			ns = ns[maxSize:]
		} else {
			ps = nil
		}
		r, err = s.upload(paths, names, u)
		if err != nil {
			continue
		}
		ss, err := u.format(r)
		if err != nil {
			continue
		}
		as = append(as, ss...)
	}
	return as, err
}

// UploadMultiAudios upload multiple VK audio attachments. It may call multiple
// times of the API. ps is arary of path string and ns is the array of the file
// name respectively to ps. ns is optional and can be nil. ns must all end with
// dot extension format.
func (s *Session) UploadMultiAudios(ps, ns []string) ([]string, error) {
	if ps == nil {
		return nil, nil
	}
	return multiUploads(s, ps, ns, getAudioUploader(0))
}

// UploadMultiAudios upload multiple VK audio attachments. It may call multiple
// times of the API. ps is array of path string and ns is the array of the file
// name respectively to ps. ns is optional and can be nil. ns must all end with
// dot extension format.
func (s *Session) UploadMultiDocs(ps, ns []string, gid int) ([]string, error) {
	if ps == nil {
		return nil, nil
	}
	return multiUploads(s, ps, ns, getDocUploader(docUploadServer, gid))
}

func (s *Session) UploadMultiWallDocs(ps, ns []string, owner int) ([]string, error) {
	if ps == nil {
		return nil, nil
	}
	return multiUploads(s, ps, ns, getDocUploader(wallDocUploadServer, owner))
}

// UploadMultiVideos upload multiple VK video attachments. It may call multiple
// times of the API. ps is arary of path string and ns is the array of the file
// name respectively to ps. ns is optional and can be nil. ns must all end with
// dot extension format.
func (s *Session) UploadMultiVideos(ps, ns []string, owner int) ([]string, error) {
	if ps == nil {
		return nil, nil
	}
	v := url.Values{}
	v.Set("group_id", strconv.Itoa(owner))
	return multiUploads(s, ps, ns, newVideoUploader(v, owner))
}

// UploadMultiPhotos upload multiple VK photo attachments. It may call multiple
// times of the API. ps is arary of path string and ns is the array of the file
// name respectively to ps. ns is optional and can be nil. These photos are
// meant for use in wall though comment can be used too. ns must all end with
// dot extension format.
func (s *Session) UploadMultiPhotos(ps, ns []string, owner int) ([]string, error) {
	if ps == nil {
		return nil, nil
	}
	return multiUploads(s, ps, ns, getWallPhotoUploader(owner))
}

func (s *Session) UploadMultiAlbumPhotos(ps, ns []string, o, a int) ([]string, error) {
	if ps == nil {
		return nil, nil
	}
	return multiUploads(s, ps, ns, getAlbumPhotoUploader(o, a))
}

// UploadMultiPMPhotos upload multiple VK photo attachments. It may call multiple
// times of the API. ps is arary of path string and ns is the array of the file
// name respectively to ps. ns is optional and can be nil. These photos are
// meant for use in private message. ns must all end with dot extension format.
func (s *Session) UploadMultiPMPhotos(ps, ns []string, gid int) ([]string, error) {
	if ps == nil {
		return nil, nil
	}
	return multiUploads(s, ps, nil, getPMPhotoUploader(gid))
}

type AttachmentUploader interface {
	AddVideo(path, name string) AttachmentUploader
	AddPhoto(path, name string) AttachmentUploader
	AddPhotos([]string) AttachmentUploader
	AddAudio(path, name string) AttachmentUploader
	AddAudios([]string) AttachmentUploader
	AddDoc(path, name string) AttachmentUploader
	AddDocs([]string) AttachmentUploader
	Upload() (string, error)
}

// This is optimized for wall posting (comments too) attachments. Video and doc
// will be uploaded to user's video section and doc section to avoid appear in
// community video and doc section.
func NewAttachmentsUploader(sess *Session, gid int) AttachmentUploader {
	return &attachUp{sess: sess, gid: gid}
}

type attachUp struct {
	gid                           int
	sess                          *Session
	p, pn, v, vn, a, an, d, dn, s []string
}

// AddVideo add a video to AttachmentUploader with p is the file path and n is
// the file name which must in the dot extension or leave n empty string.
func (a *attachUp) AddVideo(p, n string) AttachmentUploader {
	a.v = append(a.v, p)
	if n == "" {
		n = filepath.Base(p)
	}
	a.vn = append(a.vn, n)
	return a
}

// AddPhoto add a photo to AttachmentUploader with p is the file path and n is
// the file name which must in the dot extension or leave n empty string.
func (a *attachUp) AddPhoto(p, n string) AttachmentUploader {
	a.p = append(a.p, p)
	if n == "" {
		n = filepath.Base(p)
	}
	a.pn = append(a.pn, n)
	return a
}

func (a *attachUp) AddPhotos(ps []string) AttachmentUploader {
	a.p = append(a.p, ps...)
	ns := make([]string, len(ps))
	for i, p := range ps {
		ns[i] = filepath.Base(p)
	}
	a.pn = append(a.pn, ns...)
	return a
}

// AddAudio add a audio to AttachmentUploader with p is the file path and n is
// the file name which must in the dot extension or leave n empty string.
func (a *attachUp) AddAudio(p, n string) AttachmentUploader {
	a.a = append(a.a, p)
	if n == "" {
		n = filepath.Base(p)
	}
	a.an = append(a.an, n)
	return a
}

func (a *attachUp) AddAudios(ps []string) AttachmentUploader {
	a.a = append(a.a, ps...)
	ns := make([]string, len(ps))
	for i, p := range ps {
		ns[i] = filepath.Base(p)
	}
	a.an = append(a.an, ns...)
	return a
}

// AddDoc add a doc to AttachmentUploader with p is the file path and n is
// the file name which must in the dot extension or leave n empty string.
func (a *attachUp) AddDoc(p, n string) AttachmentUploader {
	a.d = append(a.d, p)
	if n == "" {
		n = filepath.Base(p)
	}
	a.dn = append(a.dn, n)
	return a
}

func (a *attachUp) AddDocs(ps []string) AttachmentUploader {
	a.d = append(a.d, ps...)
	ns := make([]string, len(ps))
	for i, p := range ps {
		ns[i] = filepath.Base(p)
	}
	a.dn = append(a.dn, ns...)
	return a
}

func multiUploadNonPhoto(a *attachUp, s *Session, gid int) error {
	var ss []string
	var err error
	if gid == 0 {
		// upload to user's
		ss, err = s.UploadMultiAudios(a.a, a.an)
		if err != nil {
			return err
		}
		a.s = append(a.s, ss...)
	} else {
		// force audio upload as doc
		a.d = append(a.d, a.a...)
		a.dn = append(a.dn, a.an...)
	}
	ss, err = s.UploadMultiVideos(a.v, a.vn, gid)
	if err != nil {
		return err
	}
	a.s = append(a.s, ss...)
	ss, err = s.UploadMultiWallDocs(a.d, a.dn, gid)
	if err != nil {
		return err
	}
	a.s = append(a.s, ss...)
	return nil
}

// Upload will generate attachment string from the uploaded files. Arg s should
// be user token.
func (a *attachUp) Upload() (string, error) {
	err := multiUploadNonPhoto(a, a.sess, 0)
	if err != nil {
		return "", err
	}
	ss, err := a.sess.UploadMultiPhotos(a.p, a.pn, 0)
	if err != nil {
		return "", err
	}
	a.s = append(a.s, ss...)
	return strings.Join(a.s, ","), nil
}

type attachPM struct {
	group *Session
	*attachUp
}

func (a *attachPM) AddVideo(p, n string) AttachmentUploader {
	a.attachUp.AddVideo(p, n)
	return a
}

func (a *attachPM) AddAudio(p, n string) AttachmentUploader {
	a.attachUp.AddDoc(p, n)
	return a
}

func (a *attachPM) AddAudios(ps []string) AttachmentUploader {
	a.attachUp.AddDocs(ps)
	return a
}

func (a *attachPM) AddDoc(p, n string) AttachmentUploader {
	a.attachUp.AddDoc(p, n)
	return a
}

func (a *attachPM) AddDocs(ps []string) AttachmentUploader {
	a.attachUp.AddDocs(ps)
	return a
}

func (a *attachPM) AddPhoto(p, n string) AttachmentUploader {
	a.attachUp.AddPhoto(p, n)
	return a
}

func (a *attachPM) AddPhotos(ps []string) AttachmentUploader {
	a.attachUp.AddPhotos(ps)
	return a
}

func (a *attachPM) Upload() (string, error) {
	gid := a.attachUp.gid
	userSess := a.attachUp.sess
	ss, err := userSess.UploadMultiVideos(a.v, a.vn, gid)
	if err != nil {
		return "", err
	}
	a.s = append(a.s, ss...)
	ss, err = userSess.UploadMultiWallDocs(a.d, a.dn, gid)
	if err != nil {
		return "", err
	}
	a.s = append(a.s, ss...)
	ss, err = a.group.UploadMultiPMPhotos(a.p, a.pn, gid)
	if err != nil {
		return "", err
	}
	a.s = append(a.s, ss...)
	return strings.Join(a.s, ","), nil
}

// NewCommunityPMUploader optimize for community PM file upload. Audio is not
// supported and it will be uploaded as doc. Only file with community ownership
// can apear on the community PM. Video upload will put the video into video
// section of the respective group. Both user and group token are required as
// group token is used for upload photo through message upload server. There
// isn't any message upload server for doc or video. Doc and video will use
// user token to upload with group_id provide to the http client params. Photo
// message upload server do not allow group_id param. Uploaded doc will also
// appear on group doc section. Use unlinker to unlink them after message with
// attachments have been sent out.
func NewCommunityPMUploader(user, group *Session, gid int) AttachmentUploader {
	return &attachPM{
		attachUp: &attachUp{gid: gid, sess: user},
		group:    group,
	}
}

const (
	UnlinkPhoto = iota
	UnlinkDoc
	UnlinkVideo
	UnlinkAudio
)

type unlinker interface {
	unlink(s string) error
}

type unlinkAttachment struct {
	s, p []string
	sess *Session
	v    url.Values
}

func NewUnlinker(sess *Session, ts ...int) unlinker {
	tsSize := len(ts)
	if tsSize <= 0 {
		return nil
	}
	ss := make([]string, tsSize)
	ps := make([]string, tsSize)
	var used int
	for _, t := range ts {
		switch t {
		case UnlinkPhoto:
			ss[used] = "photo"
			ps[used] = "photos"
			used++
		case UnlinkDoc:
			ss[used] = "doc"
			ps[used] = "docs"
			used++
		case UnlinkVideo:
			ss[used] = "video"
			ps[used] = "video"
			used++
		case UnlinkAudio:
			ss[used] = "audio"
			ps[used] = "audio"
			used++
		}
	}
	if used == 0 {
		return nil
	}
	return &unlinkAttachment{s: ss[:used], p: ps[:used], sess: sess, v: url.Values{}}
}

func (u *unlinkAttachment) unlink(s string) error {
	for i, p := range u.s {
		if strings.HasPrefix(s, p) {
			s = s[len(p):]
			as := strings.Split(s, "_")
			if len(as) == 2 {
				u.v.Set("owner_id", as[0])
				u.v.Set(fmt.Sprint(u.s[i], "_id"), as[1])
				err := u.sess.CallAPI(fmt.Sprint(u.p[i], ".delete"), u.v, new(Bool))
				if err != nil {
					return err
				}
				return nil
			}
		}
	}
	return nil
}

func UnlinkAttachments(u unlinker, s string) (err error) {
	as := strings.Split(s, ",")
	for _, a := range as {
		err = u.unlink(a)
	}
	return
}
