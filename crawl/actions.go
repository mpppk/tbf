package crawl

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
)

func Texts(sel interface{}, texts *[]string, opts ...chromedp.QueryOption) chromedp.Action {
	if texts == nil {
		panic("text cannot be nil")
	}

	return chromedp.QueryAfter(sel, func(ctxt context.Context, h *chromedp.TargetHandler, nodes ...*cdp.Node) error {
		if len(nodes) < 1 {
			return fmt.Errorf("selector `%s` did not return any nodes", sel)
		}

		for _, node := range nodes {
			t := ""
			for _, c := range node.Children {
				if c.NodeType == cdp.NodeTypeText {
					t += c.NodeValue
				}
			}
			if t != "" {
				*texts = append(*texts, t)
			}
		}
		return nil
	}, opts...)
}

func AttributeValueAll(sel interface{}, name string, values *[]string, ok *bool, opts ...chromedp.QueryOption) chromedp.Action {
	if values == nil {
		panic("values cannot be nil")
	}

	return chromedp.QueryAfter(sel, func(ctxt context.Context, h *chromedp.TargetHandler, nodes ...*cdp.Node) error {
		if len(nodes) < 1 {
			return errors.New("expected at least one element")
		}

		getAttributeValueFromNode := func(node *cdp.Node) (string, bool) {
			node.RLock()
			defer node.RUnlock()

			attrs := node.Attributes
			for i := 0; i < len(attrs); i += 2 {
				if attrs[i] == name {
					return attrs[i+1], true
				}
			}
			return "", false
		}

		if ok != nil {
			*ok = true
		}

		for _, node := range nodes {
			value, ok2 := getAttributeValueFromNode(node)
			if !ok2 {
				if ok != nil {
					*ok = false
				}
				continue
			}
			*values = append(*values, value)
		}

		return nil
	}, opts...)
}
