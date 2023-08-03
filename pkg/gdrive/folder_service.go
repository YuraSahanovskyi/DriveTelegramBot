package gdrive

import (
	"errors"

	"google.golang.org/api/drive/v3"
)

const query = "name = 'gdivebot' and mimeType = 'application/vnd.google-apps.folder'"

func getDirID(srv *drive.Service) (string, error) {
	dirID, err := getDir(srv)
	if err != nil {
		return createDir(srv)
	}
	return dirID, nil
}

func createDir(srv *drive.Service) (string, error) {
	folderFile := &drive.File{
		Name:     "gdivebot",
		MimeType: "application/vnd.google-apps.folder",
	}
	// Create the folder
	folder, err := srv.Files.Create(folderFile).Do()
	if err != nil {
		return "", err
	}

	return folder.Id, nil
}

func getDir(srv *drive.Service) (string, error) {
	// Search for the folder by its name
	list, err := srv.Files.List().Q(query).Do()
	if err != nil {
		return "", err
	}

	// If the folder already exists, return its ID
	if len(list.Files) > 0 {
		return list.Files[0].Id, nil
	}
	return "", errors.New("dir does not exist")
}
