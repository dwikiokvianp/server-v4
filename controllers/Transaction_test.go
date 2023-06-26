package controllers

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestGetTodayV2Transaction(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	var tests []struct {
		name string
		args args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetTodayV2Transaction(tt.args.c)
		})
	}
}
