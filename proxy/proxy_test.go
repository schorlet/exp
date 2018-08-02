package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

func withServer(fn func(string, string)) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	))
	defer server.Close()

	targetURL, err := url.Parse(server.URL)
	if err != nil {
		log.Fatal(err)
	}

	proxy := NewProxy(targetURL)
	go func() {
		if err := proxy.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()
	defer proxy.Close()

	fn(server.URL, proxy.URL)
}

func TestSimple(t *testing.T) {
	withServer(func(_, proxyURL string) {
		resp, err := http.Get(proxyURL)
		if err != nil {
			t.Fatal(err)
		}

		_, ok := resp.Header["Access-Control-Allow-Origin"]
		if !ok {
			t.Fatal("CORS header not set")
		}
	})
}

func withCDP(fn func(context.Context, *chromedp.CDP)) {
	// dataDir, err := ioutil.TempDir("", "proxy_test")
	// if err != nil {
	// log.Fatal(err)
	// }
	// defer os.RemoveAll(dataDir)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	c, err := chromedp.New(
		ctx,
		chromedp.WithRunnerOptions(
			runner.Flag("headless", true),
			runner.Flag("disable-gpu", true),
			// runner.UserDataDir(dataDir),
		),
		// chromedp.WithLog(log.Printf),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := c.Shutdown(ctx); err != nil {
			log.Print(err)
		}
		_ = c.Wait()
	}()

	fn(ctx, c)
}

func TestCDP(t *testing.T) {
	withServer(func(targetURL, proxyURL string) {
		withCDP(func(ctx context.Context, c *chromedp.CDP) {
			var res []byte
			await := func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
				return p.WithAwaitPromise(true)
			}

			err := c.Run(ctx, chromedp.Tasks{
				chromedp.Navigate(targetURL),
				chromedp.Sleep(100 * time.Millisecond),
				chromedp.Evaluate(`fetch("`+proxyURL+`");`, &res, await),
			})

			if err != nil {
				t.Fatal(err)
			}
		})
	})
}
