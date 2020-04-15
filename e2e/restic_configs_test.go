package tests

import (
	"bytes"
	"context"
	"github.com/kaiserkarel/drfs"
	"github.com/kaiserkarel/drfs/drive"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"testing"
)

func TestResticConfigs(t *testing.T)  {
	files := []string{"1n_2U2XPi2lGDaP04AXXNAW4r-T7_yM0A",
	"1LWxwRm-wfJN1Qo7YCVs0B-tFKg-tNmbl"}

	service, err := drive.NewService(context.Background())
	if err != nil {
		t.Fatalf("unable to init drive service: %s", err)
	}

	client, err := service.Take(context.Background(), 1)
	require.NoError(t, err)

	for _, file := range files {
		file, err := client.FilesService().Get(file).Fields("*").Do()
		require.NoError(t, err)
		log.Println(file)

		dFile, err := drfs.OpenCtx(context.Background(), file, service)
		require.NoError(t, err)

		log.Println(file.Name)

		var buf = bytes.Buffer{}
		_, err = io.Copy(&buf, dFile)
		require.NoError(t, err)
		log.Println(buf)
	}
}
