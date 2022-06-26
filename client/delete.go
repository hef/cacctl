package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// Delete
// from cloudprodelete2(ID, CID, SVN) in sitefunctions.js
func (c *Client) Delete(ctx context.Context, sid, cid int64, svn string, reserve bool) error {

	q := url.Values{}
	q.Set("cid", strconv.FormatInt(cid, 10))
	q.Set("sid", strconv.FormatInt(sid, 10)) // server id
	q.Set("svn", svn)                        // server name
	if reserve {                             // reserve ip address
		q.Set("reserve", "true")
	} else {
		q.Set("reserve", "false")
	}

	u := url.URL{
		Scheme: "https",
		Host:   "panel.cloudatcost.com",
		Path:   "/panel/_config/serverdeletecloudpro.php",
		//RawQuery: fmt.Sprintf("sid=%d&vmname=%s&cycle=%d", sid, vmName, cycle.int()), // order matters
		RawQuery: q.Encode(),
	}

	resp, err := c.sendRequest(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if len(body) > 0 {
			return fmt.Errorf("delete response: %s", string(body))
		}
	} else {
		io.Copy(os.Stdout, resp.Body)
		return fmt.Errorf("delete error: (%d) %s", resp.StatusCode, resp.Status)
	}

	return nil
}
