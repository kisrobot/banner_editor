package banner_editor

import (
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor/resource"
	"github.com/qor/serializable_meta"
)

var (
	registeredElements []*Element
)

type BannerEditorConfig struct {
	Elements        []string
	SettingResource *admin.Resource
}

type QorBannerEditorSettingInterface interface {
	serializable_meta.SerializableMetaInterface
}

type QorBannerEditorSetting struct {
	gorm.Model
	serializable_meta.SerializableMeta
}

type Element struct {
	Name     string
	Template string
	Resource *admin.Resource
	Context  func(context *admin.Context, setting interface{}) interface{}
}

func init() {
	admin.RegisterViewPath("github.com/qor/banner_editor/views")
}

func RegisterElement(e *Element) {
	registeredElements = append(registeredElements, e)
}

func (config *BannerEditorConfig) ConfigureQorMeta(metaor resource.Metaor) {
	if meta, ok := metaor.(*admin.Meta); ok {
		meta.Type = "banner_editor"
		Admin := meta.GetBaseResource().(*admin.Resource).GetAdmin()

		if config.SettingResource == nil {
			config.SettingResource = Admin.NewResource(&QorBannerEditorSetting{})
		}

		router := Admin.GetRouter()
		res := config.SettingResource
		Admin.RegisterResourceRouters(res, "read")
		router.Get(fmt.Sprintf("%v/new", res.ToParam()), New, &admin.RouteConfig{Resource: res})
		router.Post(fmt.Sprintf("%v", res.ToParam()), Create, &admin.RouteConfig{Resource: res})

		Admin.RegisterFuncMap("banner_editor_configure", func() string {
			type element struct {
				Name      string
				CreateUrl string
			}
			elements := []element{}
			newElementURL := router.Prefix + fmt.Sprintf("/%v/new", res.ToParam())
			for _, e := range registeredElements {
				elements = append(elements, element{Name: e.Name, CreateUrl: fmt.Sprintf("%v?kind=%v", newElementURL, template.URLQueryEscaper(e.Name))})
			}
			results, err := json.Marshal(elements)
			if err != nil {
				return err.Error()
			}
			return string(results)
		})
	}
}

func GetElement(name string) *Element {
	for _, e := range registeredElements {
		if e.Name == name {
			return e
		}
	}
	return nil
}

func (setting QorBannerEditorSetting) GetSerializableArgumentResource() *admin.Resource {
	element := GetElement(setting.Kind)
	if element != nil {
		return element.Resource
	}
	return nil
}
