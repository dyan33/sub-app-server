package proxy

import (
	"fmt"
	"net/url"
	"testing"
)

func TestQueryParse(t *testing.T) {

	values, err := url.ParseQuery("a=aa&b=bb")

	fmt.Println(values, err)

}
