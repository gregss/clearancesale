package excel

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"sale/internal/app"
)

func ReadByte(b []byte) ([]*app.Params, error) {
	reader := bytes.NewReader(b)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err) // todo добавлять в ошибки?
		}
	}()

	return makeParams(f)
}

func ReadFile(fileName string) ([]*app.Params, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err) // todo добавлять в ошибки?
		}
	}()

	return makeParams(f)
}

func makeParams(f *excelize.File) ([]*app.Params, error) { //nolint:gocognit
	rows, err := f.GetRows("Лист1")
	if err != nil {
		return nil, err
	}

	paramsSets := make([]*app.Params, 0)
	for _, row := range rows {
		iserr := false
		paramsSet := &app.Params{}
		for i, cell := range row {
			switch i {
			case 0:
				u, err := uuid.Parse(cell)
				if err != nil {
					// todo log.warning()
					iserr = true
					break
				}
				paramsSet.NomenclatureUUID = u.String()
			case 1:
				t, err := strconv.ParseUint(cell, 10, 32)
				if err != nil {
					// todo log.warning()
					iserr = true
					break
				}
				paramsSet.Price = uint(t)
			case 2:
				t, err := strconv.ParseUint(cell, 10, 32)
				if err != nil {
					// todo log.warning()
					iserr = true
					break
				}
				paramsSet.MaxCount = uint(t)
			case 3:
				t, err := strconv.ParseUint(cell, 10, 32)
				if err != nil {
					// todo log.warning()
					iserr = true
					break
				}
				paramsSet.MaxOrderCount = uint(t)
			case 4:
				t, err := strconv.ParseUint(cell, 10, 32)
				if err != nil {
					// todo log.warning()
					iserr = true
					break
				}
				paramsSet.Type = uint(t)
			case 5:
				t, err := strconv.ParseBool(cell)
				if err != nil {
					// todo log.warning()
					iserr = true
					break
				}
				paramsSet.IsFeed = t
			default:
				// todo log.warning()
			}
		}
		if iserr {
			continue
		}

		paramsSets = append(paramsSets, paramsSet)
	}

	return paramsSets, nil
}
