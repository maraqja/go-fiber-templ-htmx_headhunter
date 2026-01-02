package sitemap

import (
	"bytes"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sabloger/sitemap-generator/smg"
)

type SiteMapHandler struct {
	Router fiber.Router
}

func NewHandler(router fiber.Router) *SiteMapHandler {
	h := &SiteMapHandler{
		Router: router,
	}

	h.Router.Get("/sitemap.xml", h.generateSitemap)

	return h
}

func (h *SiteMapHandler) generateSitemap(c *fiber.Ctx) error {
	now := time.Now().UTC()

	sm := smg.NewSitemap(false)
	sm.SetHostname("https://example.com")
	sm.SetLastMod(&now)
	sm.SetCompress(false)
	sm.Add(&smg.SitemapLoc{
		Loc:        "/",
		LastMod:    &now,
		ChangeFreq: smg.Daily,
		Priority:   0.8,
	})
	sm.Add(&smg.SitemapLoc{
		Loc:        "/login",
		LastMod:    &now,
		ChangeFreq: smg.Weekly,
		Priority:   1,
	})

	sm.Finalize()
	var buf bytes.Buffer
	n, err := sm.WriteTo(&buf)
	if err != nil {
		return err
	}
	c.Response().Header.Set("Content-Type", "application/xml")
	c.Response().Header.Set("Content-Length", strconv.FormatInt(n, 10))
	return c.Send(buf.Bytes())
}
