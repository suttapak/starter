package repository

import (
	"bytes"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/suttapak/starter/config"
	"github.com/suttapak/starter/errs"
)

type (
	odt struct {
		conf   *config.Config
		client *resty.Client
	}

	Odt interface {
		UploadTemplate(*bytes.Reader) (id string, err error)
		UpdateTemplate(id string, r *bytes.Reader) error
		RenderPDF(id string, data map[string]interface{}, filename ...string) (res *RenderPdfData, err error)
	}

	RenderPdfData struct {
		Data    string `json:"data"`
		Message string `json:"message"`
	}

	UploadPdfData struct {
		Data struct {
			Acknowledged bool   `json:"acknowledged"`
			InsertedID   string `json:"insertedId"`
		} `json:"data"`
		Message string `json:"message"`
	}
)

// RenderPDF implements LabODTRepository.
func (l odt) RenderPDF(id string, data map[string]interface{}, filename ...string) (result *RenderPdfData, err error) {
	path, err := filepath.Abs("./public/static/pdf")
	if err != nil {
		return nil, errs.New(500, "Failed to get absolute path: "+err.Error())
	}
	fileName := ""
	if len(filename) > 0 && filename[0] != "" {
		fileName = filename[0]
	}
	result = &RenderPdfData{}
	url := l.conf.LabODT.Url + "/report/render/" + id
	res, err := l.client.R().
		SetResult(&result).
		SetAuthScheme("Bearer").
		SetHeader("Accept", "application/json").
		SetAuthToken(l.conf.LabODT.Token).
		SetBody(data).
		Post(url)
	if err != nil {
		return nil, err
	}
	if !res.IsSuccess() {
		return nil, errs.New(res.StatusCode(), string(res.Body()))
	}
	if fileName == "" {
		fileName = filepath.Base(result.Data)
	}
	res, err = l.client.SetOutputDirectory(path).R().
		SetOutput(filepath.Join(path, fileName)).
		SetAuthScheme("Bearer").
		SetAuthToken(l.conf.LabODT.Token).
		Get(l.conf.LabODT.Url + result.Data)

	if !res.IsSuccess() {
		return nil, errs.New(res.StatusCode(), string(res.Body()))
	}
	return &RenderPdfData{
		Data:    "/static/pdf/" + fileName,
		Message: result.Message,
	}, nil
}

// UpdateTemplate implements LabODTRepository.
func (l odt) UpdateTemplate(id string, r *bytes.Reader) error {

	url := l.conf.LabODT.Url + "/report/update/" + id
	res, err := l.client.R().
		SetAuthScheme("Bearer").
		SetAuthToken(l.conf.LabODT.Token).
		SetFileReader("file", "template.odt", r).
		Post(url)
	if res.StatusCode()-200 > 100 {
		return errs.New(res.StatusCode(), string(res.Body()))
	}
	return err

}

// UploadTemplate implements LabODTRepository.
func (l odt) UploadTemplate(file *bytes.Reader) (id string, err error) {

	url := l.conf.LabODT.Url + "/report/upload"
	result := UploadPdfData{}
	res, err := l.client.SetCloseConnection(true).R().
		SetResult(&result).
		SetAuthScheme("Bearer").
		SetAuthToken(l.conf.LabODT.Token).
		SetFileReader("file", "template.odt", file).
		Post(url)

	if err != nil {
		return "", err
	}
	if res.StatusCode()-200 > 100 {
		return "", errs.New(res.StatusCode(), string(res.Body()))
	}

	return result.Data.InsertedID, nil
}

func NewODT(conf *config.Config) Odt {
	return odt{
		conf:   conf,
		client: resty.New(),
	}
}
