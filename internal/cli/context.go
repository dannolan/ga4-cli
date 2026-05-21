package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dannolan/ga4-cli/internal/config"
	"github.com/dannolan/ga4-cli/internal/ga4"
	analyticsadmin "google.golang.org/api/analyticsadmin/v1beta"
	"google.golang.org/api/option"
)

type appContext struct {
	ConfigDir string
	Property  string
	Format    string
	JSON      bool
	Pretty    bool
	Out       io.Writer
	store     *config.Store
}

func (a *appContext) Store() (*config.Store, error) {
	if a.store != nil {
		return a.store, nil
	}
	store, err := config.NewStore(a.ConfigDir)
	if err != nil {
		return nil, err
	}
	a.store = store
	return store, nil
}

func (a *appContext) Config() (config.Config, error) {
	store, err := a.Store()
	if err != nil {
		return config.Config{}, err
	}
	cfg, err := store.Load()
	if err != nil {
		return cfg, err
	}
	if strings.TrimSpace(a.Property) != "" {
		cfg.PropertyID = strings.TrimSpace(a.Property)
	}
	return cfg, nil
}

func (a *appContext) Client(cmdCtx context.Context) (*ga4.Client, error) {
	httpClient, err := a.HTTPClient(cmdCtx)
	if err != nil {
		return nil, err
	}
	return ga4.NewClient(cmdCtx, option.WithHTTPClient(httpClient))
}

func (a *appContext) AdminService(cmdCtx context.Context) (*analyticsadmin.Service, error) {
	httpClient, err := a.HTTPClient(cmdCtx)
	if err != nil {
		return nil, err
	}
	return analyticsadmin.NewService(cmdCtx, option.WithHTTPClient(httpClient))
}

func (a *appContext) HTTPClient(cmdCtx context.Context) (*http.Client, error) {
	store, err := a.Store()
	if err != nil {
		return nil, err
	}
	cfg, err := a.Config()
	if err != nil {
		return nil, err
	}
	legacyPath, err := config.LegacyTokenPath()
	if err != nil {
		return nil, err
	}
	return ga4.AuthenticatedHTTPClient(cmdCtx, cfg, ga4.TokenStore{
		PrimaryPath: store.TokenPath(),
		LegacyPath:  legacyPath,
	})
}

func (a *appContext) DefaultProperty() (string, error) {
	cfg, err := a.Config()
	if err != nil {
		return "", err
	}
	if cfg.PropertyID == "" {
		return "", fmt.Errorf("property ID is required; pass --property or set GA4_PROPERTY_ID")
	}
	return cfg.PropertyID, nil
}

func (a *appContext) Print(value any) error {
	out := a.Out
	if out == nil {
		out = os.Stdout
	}
	indent := ""
	if a.Pretty || !a.JSON {
		indent = "  "
	}
	var data []byte
	var err error
	if indent == "" {
		data, err = json.Marshal(value)
	} else {
		data, err = json.MarshalIndent(value, "", indent)
	}
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, string(data))
	return err
}
