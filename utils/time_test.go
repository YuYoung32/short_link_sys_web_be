/**
 * Created by YuYoung on 2023/4/4
 * Description:
 */

package utils

import "testing"

func TestGetTodayBeginTime(t *testing.T) {
	t.Log(GetTodayBeginTime())
	t.Log(GetTodayBeginTime() - 24*60*60)
}
