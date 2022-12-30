package mob

import (
	"code.olapie.com/sugar/v2/xtype"
)

type Image xtype.Image

func NewImage() *Image {
	return new(Image)
}

type Video xtype.Video

func NewVideo() *Video {
	return new(Video)
}

func (v *Video) GetImage() *Image {
	return (*Image)(v.Image)
}

type Audio xtype.Audio

func NewAudio() *Audio {
	return new(Audio)
}

type File xtype.File

func NewFile() *File {
	return new(File)
}

type WebPage xtype.WebPage

func NewWebPage() *WebPage {
	return new(WebPage)
}

type Any xtype.Any

func (a *Any) TypeName() string {
	return (*xtype.Any)(a).TypeName()
}

func (a *Any) SetImage(i *Image) {
	(*xtype.Any)(a).SetValue((*xtype.Image)(i))
}

func (a *Any) SetAudio(au *Audio) {
	(*xtype.Any)(a).SetValue((*xtype.Audio)(au))
}

func (a *Any) SetVideo(v *Video) {
	(*xtype.Any)(a).SetValue((*xtype.Video)(v))
}

func (a *Any) SetFile(f *File) {
	(*xtype.Any)(a).SetValue((*xtype.File)(f))
}

func (a *Any) SetWebPage(wp *WebPage) {
	(*xtype.Any)(a).SetValue((*xtype.WebPage)(wp))
}

func (a *Any) Image() *Image {
	return (*Image)((*xtype.Any)(a).Value().(*xtype.Image))
}

func (a *Any) Video() *Video {
	return (*Video)((*xtype.Any)(a).Value().(*xtype.Video))
}

func (a *Any) Audio() *Audio {
	return (*Audio)((*xtype.Any)(a).Value().(*xtype.Audio))
}

func (a *Any) File() *File {
	return (*File)((*xtype.Any)(a).Value().(*xtype.File))
}

func (a *Any) WebPage() *WebPage {
	return (*WebPage)((*xtype.Any)(a).Value().(*xtype.WebPage))
}

func NewAnyObj() *Any {
	return new(Any)
}
