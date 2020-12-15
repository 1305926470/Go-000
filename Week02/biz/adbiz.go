package biz

import (
	"Week02/dao"
	"github.com/pkg/errors"
)

// default advertisement
var defaultAd = map[string]interface{}{
	"id":   1,
	"name": "php official website",
	"link": "https://www.php.net/",
	"desc": "PHP is the best language in the world",
}

// Get advertisement information based on the advertisement ID.
// return default advertisement if not found
func Advertisement(id int) (map[string]interface{}, error) {
	ret := make(map[string]interface{}, 1)
	ad, err := dao.GetAd(id)

	if err != nil {
		// return default advertisement
		if errors.Is(err, dao.ErrNoRows) {
			return defaultAd, nil
		}
		return nil, err
	}

	ret = map[string]interface{}{
		"id":   ad.Id,
		"name": ad.Name,
		"link": ad.Link,
		"desc": ad.Desc,
	}
	return ret, nil
}
