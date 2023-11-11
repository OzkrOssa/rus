package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type SaeplusResponse struct {
	Success bool        `json:"success"`
	Data    SaeplusData `json:"data"`
}

type SaeplusData struct {
	Resultado string        `json:"resultado"`
	Info      []SaeplusUser `json:"info"`
}

type SaeplusUser struct {
	IDContrato     string `json:"id_contrato"`
	NroContrato    string `json:"nro_contrato"`
	Nombre         string `json:"nombre"`
	Apellido       string `json:"apellido"`
	Cedula         string `json:"cedula"`
	InicialDoc     string `json:"inicial_doc"`
	Direccion      string `json:"direccion"`
	Telefono       string `json:"telefono"`
	TelfCasa       string `json:"telf_casa"`
	TelfAdic       string `json:"telf_adic"`
	StatusContrato string `json:"status_contrato"`
	Suscripcion    string `json:"suscripcion"`
}

type SaeplusService struct {
	attributes []string
}

func NewSaeplusService(attributes ...string) SaeplusService {
	return SaeplusService{
		attributes: attributes,
	}
}

func (ss *SaeplusService) FetchAPI() ([]SaeplusUser, error) {
	var wg sync.WaitGroup
	var saeplusUsers []SaeplusUser

	dataChan := make(chan []SaeplusUser, len(ss.attributes))
	errorChan := make(chan error, len(ss.attributes))

	for _, attr := range ss.attributes {
		wg.Add(1)
		go func(att string) {
			defer wg.Done()
			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf(
					"%s?estatus_contrato=%s",
					os.Getenv("SAEPLUS_ENDPOINT"), att,
				),
				nil,
			)
			if err != nil {
				errorChan <- fmt.Errorf("error al realizar la solicitud: %s", err)
				return
			}

			req.Header.Set(
				os.Getenv("SAEPLUS_TOKEN_HEADER"),
				os.Getenv("SAEPLUS_TOKEN"),
			)
			req.Header.Set(os.Getenv("SAEPLUS_API_HEADER"), os.Getenv("SAEPLUS_API_CONNECT"))

			client := http.Client{}
			res, getErr := client.Do(req)
			if getErr != nil {
				errorChan <- fmt.Errorf("error al realizar la solicitud HTTP: %s", getErr)
				return
			}

			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				errorChan <- fmt.Errorf("error al leer la respuesta HTTP: %s", err)
				return
			}

			var response SaeplusResponse
			err = json.Unmarshal(body, &response)
			if err != nil {
				errorChan <- fmt.Errorf("error al deserializar el JSON: %s", err)
				return
			}
			dataChan <- response.Data.Info
		}(attr)
	}

	go func() {
		wg.Wait()
		close(dataChan)
		close(errorChan)
	}()

	for data := range dataChan {
		for i := range data {
			abonado, err := strconv.Atoi(data[i].NroContrato)
			if err != nil {
				return []SaeplusUser{}, err
			}
			data[i].NroContrato = strconv.Itoa(abonado)
		}
		saeplusUsers = append(saeplusUsers, data...)
	}

	for err := range errorChan {
		return []SaeplusUser{}, err
	}

	return saeplusUsers, nil
}
