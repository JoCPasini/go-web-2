package main

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

var Trans []Transaccion

type Transaccion struct {
	Id                int      `json:"id" `
	CodigoTransaccion string   `json:"codigoTransaccion" binding:"required"`
	Moneda            []string `json:"moneda" binding:"required"`
	Monto             float64  `json:"monto" binding:"required"`
	Emisor            string   `json:"emisor" binding:"required"`
	Receptor          string   `json:"receptor" binding:"required"`
	FechaTransaccion  string   `json:"fechaTransaccion" binding:"required"`
}

func guardar(ctx *gin.Context) {

	bandera := validarToken(ctx)
	if bandera {
		return
	}
	var req Transaccion

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var errorPers string
		v := reflect.ValueOf(req)
		tipoObtenidoDeReflection := v.Type()
		//Se rompe la condición cuando le pasamos el arreglo "Moneda[]" vacío
		for i := 1; i < v.NumField(); i++ {
			if v.Field(i).Interface() == reflect.Zero(v.Field(i).Type()).Interface() {
				errorPers += fmt.Sprintf("El campo %s es requerido", tipoObtenidoDeReflection.Field(i).Name)
			}
		}
		ctx.JSON(400, gin.H{"error": errorPers})
		return
	}
	id := buscarIdMax(Trans)
	req.Id = id
	Trans = append(Trans, req)

	ctx.JSON(200, Trans)
}

func buscarIdMax(transacciones []Transaccion) int {
	id := 0
	for _, tr := range transacciones {
		if tr.Id > id {
			id = tr.Id
		}
	}
	id++
	return id
}

func validarToken(ctx *gin.Context) bool {
	token := ctx.GetHeader("token")
	if token != "123456" {
		ctx.JSON(401, gin.H{"error": "invalid token"})
		return true
	}
	return false
}

func main() {
	router := gin.Default()
	router.POST("/guardar", guardar)
	router.Run()
}
