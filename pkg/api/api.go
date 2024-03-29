package api

import (
	"github.com/gorilla/sessions"
	"msw-open-music/pkg/commonconfig"
	"msw-open-music/pkg/database"
	"msw-open-music/pkg/tmpfs"
	"net/http"
)

type API struct {
	Db                 *database.Database
	Server             http.Server
	APIConfig          commonconfig.APIConfig
	Tmpfs              *tmpfs.Tmpfs
	store              *sessions.CookieStore
	defaultSessionName string
}

func NewAPI(config commonconfig.Config) (*API, error) {
	var err error

	apiConfig := config.APIConfig
	tmpfsConfig := config.TmpfsConfig

	db, err := database.NewDatabase(apiConfig.DatabaseName, apiConfig.SingleThread)
	if err != nil {
		return nil, err
	}

	store := sessions.NewCookieStore([]byte(config.APIConfig.SECRET))

	mux := http.NewServeMux()
	apiMux := http.NewServeMux()

	api := &API{
		Db: db,
		Server: http.Server{
			Addr:    apiConfig.Addr,
			Handler: mux,
		},
		APIConfig:          apiConfig,
		store:              store,
		defaultSessionName: "msw-open-music",
	}
	api.Tmpfs = tmpfs.NewTmpfs(tmpfsConfig)

	// mount api
	apiMux.HandleFunc("/hello", api.HandleOK)
	apiMux.HandleFunc("/get_file", api.HandleGetFile)
	apiMux.HandleFunc("/get_file_direct", api.HandleGetFileDirect)
	apiMux.HandleFunc("/search_files", api.HandleSearchFiles)
	apiMux.HandleFunc("/search_folders", api.HandleSearchFolders)
	apiMux.HandleFunc("/get_files_in_folder", api.HandleGetFilesInFolder)
	apiMux.HandleFunc("/get_random_files", api.HandleGetRandomFiles)
	apiMux.HandleFunc("/get_random_files_with_tag", api.HandleGetRandomFilesWithTag)
	apiMux.HandleFunc("/get_file_stream", api.HandleGetFileStream)
	apiMux.HandleFunc("/get_ffmpeg_config_list", api.HandleGetFfmpegConfigs)
	apiMux.HandleFunc("/get_file_info", api.HandleGetFileInfo)
	apiMux.HandleFunc("/get_file_ffprobe_info", api.HandleGetFileFfprobeInfo)
	apiMux.HandleFunc("/get_file_stream_direct", api.HandleGetFileStreamDirect)
	apiMux.HandleFunc("/prepare_file_stream_direct", api.HandlePrepareFileStreamDirect)
	apiMux.HandleFunc("/delete_file", api.HandleDeleteFile)
	apiMux.HandleFunc("/update_filename", api.HandleUpdateFilename)
	apiMux.HandleFunc("/reset_filename", api.HandleResetFilename)
	apiMux.HandleFunc("/reset_foldername", api.HandleResetFoldername)
	// feedback
	apiMux.HandleFunc("/feedback", api.HandleFeedback)
	apiMux.HandleFunc("/get_feedbacks", api.HandleGetFeedbacks)
	apiMux.HandleFunc("/delete_feedback", api.HandleDeleteFeedback)
	// user
	apiMux.HandleFunc("/login", api.HandleLogin)
	apiMux.HandleFunc("/register", api.HandleRegister)
	apiMux.HandleFunc("/logout", api.HandleLoginAsAnonymous)
	apiMux.HandleFunc("/get_user_info", api.HandleGetUserInfo)
	apiMux.HandleFunc("/get_users", api.HandleGetUsers)
	apiMux.HandleFunc("/update_user_active", api.HandleUpdateUserActive)
	apiMux.HandleFunc("/update_username", api.HandleUpdateUsername)
	apiMux.HandleFunc("/update_user_password", api.HandleUpdateUserPassword)
	// tag
	apiMux.HandleFunc("/get_tags", api.HandleGetTags)
	apiMux.HandleFunc("/get_tag_info", api.HandleGetTagInfo)
	apiMux.HandleFunc("/insert_tag", api.HandleInsertTag)
	apiMux.HandleFunc("/update_tag", api.HandleUpdateTag)
	apiMux.HandleFunc("/put_tag_on_file", api.HandlePutTagOnFile)
	apiMux.HandleFunc("/get_tags_on_file", api.HandleGetTagsOnFile)
	apiMux.HandleFunc("/delete_tag_on_file", api.HandleDeleteTagOnFile)
	apiMux.HandleFunc("/delete_tag", api.HandleDeleteTag)
	// folder
	apiMux.HandleFunc("/update_foldername", api.HandleUpdateFoldername)
	// review
	apiMux.HandleFunc("/insert_review", api.HandleInsertReview)
	apiMux.HandleFunc("/get_reviews_on_file", api.HandleGetReviewsOnFile)
	apiMux.HandleFunc("/get_review", api.HandleGetReview)
	apiMux.HandleFunc("/update_review", api.HandleUpdateReview)
	apiMux.HandleFunc("/delete_review", api.HandleDeleteReview)
	apiMux.HandleFunc("/get_reviews_by_user", api.HandleGetReviewsByUser)
	// statistic
	apiMux.HandleFunc("/record_playback", api.HandleRecordPlayback)
	// database
	apiMux.HandleFunc("/walk", api.HandleWalk)
	apiMux.HandleFunc("/reset", api.HandleReset)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", api.PermissionMiddleware(apiMux)))
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("web/build"))))

	return api, nil
}
