package utils

import (
	"context"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadToCloudinary(file []byte) (string, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return "", err
	}

	// Convert bool to *bool
	useFalse := false
	uniqueFalse := false

	uploadResult, err := cld.Upload.Upload(
		context.Background(),
		file,
		uploader.UploadParams{
			Folder:         "user_avatars",
			UseFilename:    &useFalse,
			UniqueFilename: &uniqueFalse,
		})
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}
