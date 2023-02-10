package mob

import (
	"code.olapie.com/sugar/v2/types"
)

type Image types.Image

func NewImage() *Image {
	return new(Image)
}

type Video types.Video

func NewVideo() *Video {
	return new(Video)
}

func (v *Video) GetImage() *Image {
	return (*Image)(v.Image)
}

type Audio types.Audio

func NewAudio() *Audio {
	return new(Audio)
}

type File types.File

func NewFile() *File {
	return new(File)
}

type WebPage types.WebPage

func NewWebPage() *WebPage {
	return new(WebPage)
}

type Any types.Any

func (a *Any) TypeName() string {
	return (*types.Any)(a).TypeName()
}

func (a *Any) SetImage(i *Image) {
	(*types.Any)(a).SetValue((*types.Image)(i))
}

func (a *Any) SetAudio(au *Audio) {
	(*types.Any)(a).SetValue((*types.Audio)(au))
}

func (a *Any) SetVideo(v *Video) {
	(*types.Any)(a).SetValue((*types.Video)(v))
}

func (a *Any) SetFile(f *File) {
	(*types.Any)(a).SetValue((*types.File)(f))
}

func (a *Any) SetWebPage(wp *WebPage) {
	(*types.Any)(a).SetValue((*types.WebPage)(wp))
}

func (a *Any) Image() *Image {
	return (*Image)((*types.Any)(a).Value().(*types.Image))
}

func (a *Any) Video() *Video {
	return (*Video)((*types.Any)(a).Value().(*types.Video))
}

func (a *Any) Audio() *Audio {
	return (*Audio)((*types.Any)(a).Value().(*types.Audio))
}

func (a *Any) File() *File {
	return (*File)((*types.Any)(a).Value().(*types.File))
}

func (a *Any) WebPage() *WebPage {
	return (*WebPage)((*types.Any)(a).Value().(*types.WebPage))
}

func NewAnyObj() *Any {
	return new(Any)
}
