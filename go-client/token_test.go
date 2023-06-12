/*
 *   Copyright (c) 2022 Intel Corporation
 *   All rights reserved.
 *   SPDX-License-Identifier: BSD-3-Clause
 */
package client

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"net/http"
	"testing"
	"time"
)

var (
	token             = "eyJhbGciOiJQUzM4NCIsImprdSI6Imh0dHBzOi8vd3d3LmludGVsLmNvbS9hbWJlci9jZXJ0cyIsImtpZCI6IjNjMjQxOGI1ZTY5ZTI2NDRiOTE2NzJmZjYwNTY2NjRkOTI0MjM0ZjAiLCJ0eXAiOiJKV1QifQ.eyJhbWJlcl90cnVzdF9zY29yZSI6MTAsImFtYmVyX3JlcG9ydF9kYXRhIjoiZWZmNWEyYTExNDg2N2FhOTQ0NjIwYzQ4Y2Q4NjcwNDZkYmY2ZjdmY2JmODQ5YTliNjZhNzg3MjJmNGRjZDdjOTAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLCJhbWJlcl90ZWVfaGVsZF9kYXRhIjoiQVFBQkFLMkxyUDhSZ1ROMW5naFJZWTBKQnZZQ1M3K2JCcVhoRjIzY2JkVFVRR3F3MEl2Wm9NYkpySUNmQXJ1MjVWWWpKbkZaS0Vvb1hRWmhPZUlBeTZNV0RpWEpmdDc0VnVUR25YZnNLUDk4bWNvZXBiQ2M4U1BJRFBsdkhTQy9QQWtlRzJUdlZ3QkhBbjcvcURVNXJwUENIS05xWDYweUd2SW95QjhrLzBJTzR5M2V0ekdvQjF5YVRPQ3Iyd1NCYmdUUkV1M3ppY3JJODFPL1RsK0FaWitVekFCdUxSUEgxdlBrelBQYVhCT21IN2Q4SnZXZ0RwSjhFenBNRitzakt4dXI1dkEraDNjamxDUG4yZjFjODhqZGIyNDgyUUJ2STZoanB4R2k0dWRUVVdJekdKdElKeElUbnNscThwTFpUekhnR3V2UEdYK2xkYWFKdUVPSGJBeFhDRW4zTnNGNHZvVjFSQ2I1OXBBMEI2NnZBd1RHZmFONE9pR205aGhiTk1NTnZNeGlhZmdGanJWWHpjc1BvUE5vN2hPd0dMcVJFdGUrMWkzZzlGNDBCK2hEZVV6elZhTU8zVkxHTUtEcDlUSDJqMytYSnRnU3p4dThOWlg1WEZVeGpSMlJINzV5d25vbnRNQStnaDZid1d1UUlWWWI2K0k3eHEzdWxOaUZldzZ4eWc9PSIsImFtYmVyX3NneF9tcmVuY2xhdmUiOiI4M2Y0ZTgxOTg2MWFkZWY2ZmZiMmE0ODY1ZWZlYTkzMzdiOTFlZDMwZmEzMzQ5MWIxN2YwZDVkOWU4MjA0NDEwIiwiYW1iZXJfc2d4X2lzX2RlYnVnZ2FibGUiOmZhbHNlLCJhbWJlcl9zZ3hfbXJzaWduZXIiOiI4M2Q3MTllNzdkZWFjYTE0NzBmNmJhZjYyYTRkNzc0MzAzYzg5OWRiNjkwMjBmOWM3MGVlMWRmYzA4YzdjZTllIiwiYW1iZXJfc2d4X2lzdnByb2RpZCI6MCwiYW1iZXJfc2d4X2lzdnN2biI6MCwiYW1iZXJfbWF0Y2hlZF9wb2xpY3lfaWRzIjpbImY0MzZjMzBhLTY0NGUtNGJiMi1iMzJjLTFmNWJmZjc3NTJmMiJdLCJhbWJlci1mYWl0aGZ1bC1zZXJ2aWNlLWlkcyI6WyI2YmFhNjMwMS0zMWVlLTQ1NmMtOWEzOC1lMjc5YWM3ZjZkNmEiLCJkZTU3NDU5ZC1mYjU2LTRhODgtYmU1ZC02ZjJhMzcwMzkzOWMiXSwiYW1iZXJfdGNiX3N0YXR1cyI6Ik9LIiwiYW1iZXJfZXZpZGVuY2VfdHlwZSI6IlNHWCIsImFtYmVyX3NpZ25lZF9ub25jZSI6dHJ1ZSwiYW1iZXJfY3VzdG9tX3BvbGljeSI6e30sInZlciI6IjEuMCIsImV4cCI6MTY2NDQ0NzA5MywianRpIjoiMzc5ZjBiMDctNTUzNy00YzdhLWFlNTAtODk2YTU1ZjIzNzY2IiwiaWF0IjoxNjY0NDQ2NzYzLCJpc3MiOiJBUyBBdHRlc3RhdGlvbiBUb2tlbiBJc3N1ZXIifQ.X2UDvraRVzAJpC1G1WAK2Qbx64d9WI5T_AKAq1lK5VAjEf409y5fZPxkBdZ-fGYt653nQ5Ah0-jkFRt0Yo7B2cxNmDWn61mMW9yYtt_55qHcbuDX5x4a-7MVawWjS1gLzY7qddmpzoIhwrx575c5JoQjG4qybDejRufUxhvu_XOOxSfhyh4JGRxBYNX19ZGeIbHtE3mfAXqg6qphZFfClIQLdlU-wGbefyN5mwpTK0T3eQ9Tlt0zZFrcv7lNIAPHHB52Ke9R7qdFEoVNNX-8YFMzk4gQyZdzYJS7Q3ElhQYXhBWnY5iwquftQztQcfydJL8o1OC-Ru0s-keF7OBaxABHcv5OhUKlVc44zaBnekP9lTzRCINYnVK67KxyAHsgXkK19UiX6v1FYxdcmdwZgNn5OkCwxiAMLgB8_CQku6q4aeyhCMo4acD1xKd6kkfgYQDxehLbGV6weT4E60omx6UFE13L9yANNNoWtzy0A4PJsiw0tbRPYYO8ehZ8Vgrb9sc00cdqG7_7ok-iivuxklaaSuzrY8VtkGw9T8g0w__fJ0X2KCMPcl3XfNidhOGxJ9402ff93X-QY3dHyaLOqmtJK0vlQ0vuoThseBSOezETalhFCuh-JUYZskokQ21fDPs2xDiytKubxqzrJVF1G1n1AVNWlIZXPXLyoXANS4s"
	tokenMissingKID   = "eyJhbGciOiJQUzM4NCIsImprdSI6Imh0dHBzOi8vd3d3LmludGVsLmNvbS9hbWJlci9jZXJ0cyIsInR5cCI6IkpXVCJ9.eyJhbWJlcl90cnVzdF9zY29yZSI6MTAsImFtYmVyX3JlcG9ydF9kYXRhIjoiZWZmNWEyYTExNDg2N2FhOTQ0NjIwYzQ4Y2Q4NjcwNDZkYmY2ZjdmY2JmODQ5YTliNjZhNzg3MjJmNGRjZDdjOTAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLCJhbWJlcl90ZWVfaGVsZF9kYXRhIjoiQVFBQkFLMkxyUDhSZ1ROMW5naFJZWTBKQnZZQ1M3K2JCcVhoRjIzY2JkVFVRR3F3MEl2Wm9NYkpySUNmQXJ1MjVWWWpKbkZaS0Vvb1hRWmhPZUlBeTZNV0RpWEpmdDc0VnVUR25YZnNLUDk4bWNvZXBiQ2M4U1BJRFBsdkhTQy9QQWtlRzJUdlZ3QkhBbjcvcURVNXJwUENIS05xWDYweUd2SW95QjhrLzBJTzR5M2V0ekdvQjF5YVRPQ3Iyd1NCYmdUUkV1M3ppY3JJODFPL1RsK0FaWitVekFCdUxSUEgxdlBrelBQYVhCT21IN2Q4SnZXZ0RwSjhFenBNRitzakt4dXI1dkEraDNjamxDUG4yZjFjODhqZGIyNDgyUUJ2STZoanB4R2k0dWRUVVdJekdKdElKeElUbnNscThwTFpUekhnR3V2UEdYK2xkYWFKdUVPSGJBeFhDRW4zTnNGNHZvVjFSQ2I1OXBBMEI2NnZBd1RHZmFONE9pR205aGhiTk1NTnZNeGlhZmdGanJWWHpjc1BvUE5vN2hPd0dMcVJFdGUrMWkzZzlGNDBCK2hEZVV6elZhTU8zVkxHTUtEcDlUSDJqMytYSnRnU3p4dThOWlg1WEZVeGpSMlJINzV5d25vbnRNQStnaDZid1d1UUlWWWI2K0k3eHEzdWxOaUZldzZ4eWc9PSIsImFtYmVyX3NneF9tcmVuY2xhdmUiOiI4M2Y0ZTgxOTg2MWFkZWY2ZmZiMmE0ODY1ZWZlYTkzMzdiOTFlZDMwZmEzMzQ5MWIxN2YwZDVkOWU4MjA0NDEwIiwiYW1iZXJfc2d4X2lzX2RlYnVnZ2FibGUiOmZhbHNlLCJhbWJlcl9zZ3hfbXJzaWduZXIiOiI4M2Q3MTllNzdkZWFjYTE0NzBmNmJhZjYyYTRkNzc0MzAzYzg5OWRiNjkwMjBmOWM3MGVlMWRmYzA4YzdjZTllIiwiYW1iZXJfc2d4X2lzdnByb2RpZCI6MCwiYW1iZXJfc2d4X2lzdnN2biI6MCwiYW1iZXJfbWF0Y2hlZF9wb2xpY3lfaWRzIjpbImY0MzZjMzBhLTY0NGUtNGJiMi1iMzJjLTFmNWJmZjc3NTJmMiJdLCJhbWJlci1mYWl0aGZ1bC1zZXJ2aWNlLWlkcyI6WyI2YmFhNjMwMS0zMWVlLTQ1NmMtOWEzOC1lMjc5YWM3ZjZkNmEiLCJkZTU3NDU5ZC1mYjU2LTRhODgtYmU1ZC02ZjJhMzcwMzkzOWMiXSwiYW1iZXJfdGNiX3N0YXR1cyI6Ik9LIiwiYW1iZXJfZXZpZGVuY2VfdHlwZSI6IlNHWCIsImFtYmVyX3NpZ25lZF9ub25jZSI6dHJ1ZSwiYW1iZXJfY3VzdG9tX3BvbGljeSI6e30sInZlciI6IjEuMCIsImV4cCI6MTY2NDQ0NzA5MywianRpIjoiMzc5ZjBiMDctNTUzNy00YzdhLWFlNTAtODk2YTU1ZjIzNzY2IiwiaWF0IjoxNjY0NDQ2NzYzLCJpc3MiOiJBUyBBdHRlc3RhdGlvbiBUb2tlbiBJc3N1ZXIifQ.X2UDvraRVzAJpC1G1WAK2Qbx64d9WI5T_AKAq1lK5VAjEf409y5fZPxkBdZ-fGYt653nQ5Ah0-jkFRt0Yo7B2cxNmDWn61mMW9yYtt_55qHcbuDX5x4a-7MVawWjS1gLzY7qddmpzoIhwrx575c5JoQjG4qybDejRufUxhvu_XOOxSfhyh4JGRxBYNX19ZGeIbHtE3mfAXqg6qphZFfClIQLdlU-wGbefyN5mwpTK0T3eQ9Tlt0zZFrcv7lNIAPHHB52Ke9R7qdFEoVNNX-8YFMzk4gQyZdzYJS7Q3ElhQYXhBWnY5iwquftQztQcfydJL8o1OC-Ru0s-keF7OBaxABHcv5OhUKlVc44zaBnekP9lTzRCINYnVK67KxyAHsgXkK19UiX6v1FYxdcmdwZgNn5OkCwxiAMLgB8_CQku6q4aeyhCMo4acD1xKd6kkfgYQDxehLbGV6weT4E60omx6UFE13L9yANNNoWtzy0A4PJsiw0tbRPYYO8ehZ8Vgrb9sc00cdqG7_7ok-iivuxklaaSuzrY8VtkGw9T8g0w__fJ0X2KCMPcl3XfNidhOGxJ9402ff93X-QY3dHyaLOqmtJK0vlQ0vuoThseBSOezETalhFCuh-JUYZskokQ21fDPs2xDiytKubxqzrJVF1G1n1AVNWlIZXPXLyoXANS4s"
	tokenInvalidKID   = "ewogICJhbGciOiAiUFMzODQiLAogICJqa3UiOiAiaHR0cHM6Ly93d3cuaW50ZWwuY29tL2FtYmVyL2NlcnRzIiwKICAia2lkIjogMTIzLAogICJ0eXAiOiAiSldUIgp9.eyJhbWJlcl90cnVzdF9zY29yZSI6MTAsImFtYmVyX3JlcG9ydF9kYXRhIjoiZWZmNWEyYTExNDg2N2FhOTQ0NjIwYzQ4Y2Q4NjcwNDZkYmY2ZjdmY2JmODQ5YTliNjZhNzg3MjJmNGRjZDdjOTAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLCJhbWJlcl90ZWVfaGVsZF9kYXRhIjoiQVFBQkFLMkxyUDhSZ1ROMW5naFJZWTBKQnZZQ1M3K2JCcVhoRjIzY2JkVFVRR3F3MEl2Wm9NYkpySUNmQXJ1MjVWWWpKbkZaS0Vvb1hRWmhPZUlBeTZNV0RpWEpmdDc0VnVUR25YZnNLUDk4bWNvZXBiQ2M4U1BJRFBsdkhTQy9QQWtlRzJUdlZ3QkhBbjcvcURVNXJwUENIS05xWDYweUd2SW95QjhrLzBJTzR5M2V0ekdvQjF5YVRPQ3Iyd1NCYmdUUkV1M3ppY3JJODFPL1RsK0FaWitVekFCdUxSUEgxdlBrelBQYVhCT21IN2Q4SnZXZ0RwSjhFenBNRitzakt4dXI1dkEraDNjamxDUG4yZjFjODhqZGIyNDgyUUJ2STZoanB4R2k0dWRUVVdJekdKdElKeElUbnNscThwTFpUekhnR3V2UEdYK2xkYWFKdUVPSGJBeFhDRW4zTnNGNHZvVjFSQ2I1OXBBMEI2NnZBd1RHZmFONE9pR205aGhiTk1NTnZNeGlhZmdGanJWWHpjc1BvUE5vN2hPd0dMcVJFdGUrMWkzZzlGNDBCK2hEZVV6elZhTU8zVkxHTUtEcDlUSDJqMytYSnRnU3p4dThOWlg1WEZVeGpSMlJINzV5d25vbnRNQStnaDZid1d1UUlWWWI2K0k3eHEzdWxOaUZldzZ4eWc9PSIsImFtYmVyX3NneF9tcmVuY2xhdmUiOiI4M2Y0ZTgxOTg2MWFkZWY2ZmZiMmE0ODY1ZWZlYTkzMzdiOTFlZDMwZmEzMzQ5MWIxN2YwZDVkOWU4MjA0NDEwIiwiYW1iZXJfc2d4X2lzX2RlYnVnZ2FibGUiOmZhbHNlLCJhbWJlcl9zZ3hfbXJzaWduZXIiOiI4M2Q3MTllNzdkZWFjYTE0NzBmNmJhZjYyYTRkNzc0MzAzYzg5OWRiNjkwMjBmOWM3MGVlMWRmYzA4YzdjZTllIiwiYW1iZXJfc2d4X2lzdnByb2RpZCI6MCwiYW1iZXJfc2d4X2lzdnN2biI6MCwiYW1iZXJfbWF0Y2hlZF9wb2xpY3lfaWRzIjpbImY0MzZjMzBhLTY0NGUtNGJiMi1iMzJjLTFmNWJmZjc3NTJmMiJdLCJhbWJlci1mYWl0aGZ1bC1zZXJ2aWNlLWlkcyI6WyI2YmFhNjMwMS0zMWVlLTQ1NmMtOWEzOC1lMjc5YWM3ZjZkNmEiLCJkZTU3NDU5ZC1mYjU2LTRhODgtYmU1ZC02ZjJhMzcwMzkzOWMiXSwiYW1iZXJfdGNiX3N0YXR1cyI6Ik9LIiwiYW1iZXJfZXZpZGVuY2VfdHlwZSI6IlNHWCIsImFtYmVyX3NpZ25lZF9ub25jZSI6dHJ1ZSwiYW1iZXJfY3VzdG9tX3BvbGljeSI6e30sInZlciI6IjEuMCIsImV4cCI6MTY2NDQ0NzA5MywianRpIjoiMzc5ZjBiMDctNTUzNy00YzdhLWFlNTAtODk2YTU1ZjIzNzY2IiwiaWF0IjoxNjY0NDQ2NzYzLCJpc3MiOiJBUyBBdHRlc3RhdGlvbiBUb2tlbiBJc3N1ZXIifQ.X2UDvraRVzAJpC1G1WAK2Qbx64d9WI5T_AKAq1lK5VAjEf409y5fZPxkBdZ-fGYt653nQ5Ah0-jkFRt0Yo7B2cxNmDWn61mMW9yYtt_55qHcbuDX5x4a-7MVawWjS1gLzY7qddmpzoIhwrx575c5JoQjG4qybDejRufUxhvu_XOOxSfhyh4JGRxBYNX19ZGeIbHtE3mfAXqg6qphZFfClIQLdlU-wGbefyN5mwpTK0T3eQ9Tlt0zZFrcv7lNIAPHHB52Ke9R7qdFEoVNNX-8YFMzk4gQyZdzYJS7Q3ElhQYXhBWnY5iwquftQztQcfydJL8o1OC-Ru0s-keF7OBaxABHcv5OhUKlVc44zaBnekP9lTzRCINYnVK67KxyAHsgXkK19UiX6v1FYxdcmdwZgNn5OkCwxiAMLgB8_CQku6q4aeyhCMo4acD1xKd6kkfgYQDxehLbGV6weT4E60omx6UFE13L9yANNNoWtzy0A4PJsiw0tbRPYYO8ehZ8Vgrb9sc00cdqG7_7ok-iivuxklaaSuzrY8VtkGw9T8g0w__fJ0X2KCMPcl3XfNidhOGxJ9402ff93X-QY3dHyaLOqmtJK0vlQ0vuoThseBSOezETalhFCuh-JUYZskokQ21fDPs2xDiytKubxqzrJVF1G1n1AVNWlIZXPXLyoXANS4s"
	tokenMissingJKU   = "ewogICJhbGciOiAiUFMzODQiLAogICJraWQiOiAiM2MyNDE4YjVlNjllMjY0NGI5MTY3MmZmNjA1NjY2NGQ5MjQyMzRmMCIsCiAgInR5cCI6ICJKV1QiCn0.eyJhbWJlcl90cnVzdF9zY29yZSI6MTAsImFtYmVyX3JlcG9ydF9kYXRhIjoiZWZmNWEyYTExNDg2N2FhOTQ0NjIwYzQ4Y2Q4NjcwNDZkYmY2ZjdmY2JmODQ5YTliNjZhNzg3MjJmNGRjZDdjOTAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLCJhbWJlcl90ZWVfaGVsZF9kYXRhIjoiQVFBQkFLMkxyUDhSZ1ROMW5naFJZWTBKQnZZQ1M3K2JCcVhoRjIzY2JkVFVRR3F3MEl2Wm9NYkpySUNmQXJ1MjVWWWpKbkZaS0Vvb1hRWmhPZUlBeTZNV0RpWEpmdDc0VnVUR25YZnNLUDk4bWNvZXBiQ2M4U1BJRFBsdkhTQy9QQWtlRzJUdlZ3QkhBbjcvcURVNXJwUENIS05xWDYweUd2SW95QjhrLzBJTzR5M2V0ekdvQjF5YVRPQ3Iyd1NCYmdUUkV1M3ppY3JJODFPL1RsK0FaWitVekFCdUxSUEgxdlBrelBQYVhCT21IN2Q4SnZXZ0RwSjhFenBNRitzakt4dXI1dkEraDNjamxDUG4yZjFjODhqZGIyNDgyUUJ2STZoanB4R2k0dWRUVVdJekdKdElKeElUbnNscThwTFpUekhnR3V2UEdYK2xkYWFKdUVPSGJBeFhDRW4zTnNGNHZvVjFSQ2I1OXBBMEI2NnZBd1RHZmFONE9pR205aGhiTk1NTnZNeGlhZmdGanJWWHpjc1BvUE5vN2hPd0dMcVJFdGUrMWkzZzlGNDBCK2hEZVV6elZhTU8zVkxHTUtEcDlUSDJqMytYSnRnU3p4dThOWlg1WEZVeGpSMlJINzV5d25vbnRNQStnaDZid1d1UUlWWWI2K0k3eHEzdWxOaUZldzZ4eWc9PSIsImFtYmVyX3NneF9tcmVuY2xhdmUiOiI4M2Y0ZTgxOTg2MWFkZWY2ZmZiMmE0ODY1ZWZlYTkzMzdiOTFlZDMwZmEzMzQ5MWIxN2YwZDVkOWU4MjA0NDEwIiwiYW1iZXJfc2d4X2lzX2RlYnVnZ2FibGUiOmZhbHNlLCJhbWJlcl9zZ3hfbXJzaWduZXIiOiI4M2Q3MTllNzdkZWFjYTE0NzBmNmJhZjYyYTRkNzc0MzAzYzg5OWRiNjkwMjBmOWM3MGVlMWRmYzA4YzdjZTllIiwiYW1iZXJfc2d4X2lzdnByb2RpZCI6MCwiYW1iZXJfc2d4X2lzdnN2biI6MCwiYW1iZXJfbWF0Y2hlZF9wb2xpY3lfaWRzIjpbImY0MzZjMzBhLTY0NGUtNGJiMi1iMzJjLTFmNWJmZjc3NTJmMiJdLCJhbWJlci1mYWl0aGZ1bC1zZXJ2aWNlLWlkcyI6WyI2YmFhNjMwMS0zMWVlLTQ1NmMtOWEzOC1lMjc5YWM3ZjZkNmEiLCJkZTU3NDU5ZC1mYjU2LTRhODgtYmU1ZC02ZjJhMzcwMzkzOWMiXSwiYW1iZXJfdGNiX3N0YXR1cyI6Ik9LIiwiYW1iZXJfZXZpZGVuY2VfdHlwZSI6IlNHWCIsImFtYmVyX3NpZ25lZF9ub25jZSI6dHJ1ZSwiYW1iZXJfY3VzdG9tX3BvbGljeSI6e30sInZlciI6IjEuMCIsImV4cCI6MTY2NDQ0NzA5MywianRpIjoiMzc5ZjBiMDctNTUzNy00YzdhLWFlNTAtODk2YTU1ZjIzNzY2IiwiaWF0IjoxNjY0NDQ2NzYzLCJpc3MiOiJBUyBBdHRlc3RhdGlvbiBUb2tlbiBJc3N1ZXIifQ.X2UDvraRVzAJpC1G1WAK2Qbx64d9WI5T_AKAq1lK5VAjEf409y5fZPxkBdZ-fGYt653nQ5Ah0-jkFRt0Yo7B2cxNmDWn61mMW9yYtt_55qHcbuDX5x4a-7MVawWjS1gLzY7qddmpzoIhwrx575c5JoQjG4qybDejRufUxhvu_XOOxSfhyh4JGRxBYNX19ZGeIbHtE3mfAXqg6qphZFfClIQLdlU-wGbefyN5mwpTK0T3eQ9Tlt0zZFrcv7lNIAPHHB52Ke9R7qdFEoVNNX-8YFMzk4gQyZdzYJS7Q3ElhQYXhBWnY5iwquftQztQcfydJL8o1OC-Ru0s-keF7OBaxABHcv5OhUKlVc44zaBnekP9lTzRCINYnVK67KxyAHsgXkK19UiX6v1FYxdcmdwZgNn5OkCwxiAMLgB8_CQku6q4aeyhCMo4acD1xKd6kkfgYQDxehLbGV6weT4E60omx6UFE13L9yANNNoWtzy0A4PJsiw0tbRPYYO8ehZ8Vgrb9sc00cdqG7_7ok-iivuxklaaSuzrY8VtkGw9T8g0w__fJ0X2KCMPcl3XfNidhOGxJ9402ff93X-QY3dHyaLOqmtJK0vlQ0vuoThseBSOezETalhFCuh-JUYZskokQ21fDPs2xDiytKubxqzrJVF1G1n1AVNWlIZXPXLyoXANS4s"
	tokenInvalidJKU   = "ewogICJhbGciOiAiUFMzODQiLAogICJqa3UiOiAxMjMsCiAgImtpZCI6ICIzYzI0MThiNWU2OWUyNjQ0YjkxNjcyZmY2MDU2NjY0ZDkyNDIzNGYwIiwKICAidHlwIjogIkpXVCIKfQ.eyJhbWJlcl90cnVzdF9zY29yZSI6MTAsImFtYmVyX3JlcG9ydF9kYXRhIjoiZWZmNWEyYTExNDg2N2FhOTQ0NjIwYzQ4Y2Q4NjcwNDZkYmY2ZjdmY2JmODQ5YTliNjZhNzg3MjJmNGRjZDdjOTAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLCJhbWJlcl90ZWVfaGVsZF9kYXRhIjoiQVFBQkFLMkxyUDhSZ1ROMW5naFJZWTBKQnZZQ1M3K2JCcVhoRjIzY2JkVFVRR3F3MEl2Wm9NYkpySUNmQXJ1MjVWWWpKbkZaS0Vvb1hRWmhPZUlBeTZNV0RpWEpmdDc0VnVUR25YZnNLUDk4bWNvZXBiQ2M4U1BJRFBsdkhTQy9QQWtlRzJUdlZ3QkhBbjcvcURVNXJwUENIS05xWDYweUd2SW95QjhrLzBJTzR5M2V0ekdvQjF5YVRPQ3Iyd1NCYmdUUkV1M3ppY3JJODFPL1RsK0FaWitVekFCdUxSUEgxdlBrelBQYVhCT21IN2Q4SnZXZ0RwSjhFenBNRitzakt4dXI1dkEraDNjamxDUG4yZjFjODhqZGIyNDgyUUJ2STZoanB4R2k0dWRUVVdJekdKdElKeElUbnNscThwTFpUekhnR3V2UEdYK2xkYWFKdUVPSGJBeFhDRW4zTnNGNHZvVjFSQ2I1OXBBMEI2NnZBd1RHZmFONE9pR205aGhiTk1NTnZNeGlhZmdGanJWWHpjc1BvUE5vN2hPd0dMcVJFdGUrMWkzZzlGNDBCK2hEZVV6elZhTU8zVkxHTUtEcDlUSDJqMytYSnRnU3p4dThOWlg1WEZVeGpSMlJINzV5d25vbnRNQStnaDZid1d1UUlWWWI2K0k3eHEzdWxOaUZldzZ4eWc9PSIsImFtYmVyX3NneF9tcmVuY2xhdmUiOiI4M2Y0ZTgxOTg2MWFkZWY2ZmZiMmE0ODY1ZWZlYTkzMzdiOTFlZDMwZmEzMzQ5MWIxN2YwZDVkOWU4MjA0NDEwIiwiYW1iZXJfc2d4X2lzX2RlYnVnZ2FibGUiOmZhbHNlLCJhbWJlcl9zZ3hfbXJzaWduZXIiOiI4M2Q3MTllNzdkZWFjYTE0NzBmNmJhZjYyYTRkNzc0MzAzYzg5OWRiNjkwMjBmOWM3MGVlMWRmYzA4YzdjZTllIiwiYW1iZXJfc2d4X2lzdnByb2RpZCI6MCwiYW1iZXJfc2d4X2lzdnN2biI6MCwiYW1iZXJfbWF0Y2hlZF9wb2xpY3lfaWRzIjpbImY0MzZjMzBhLTY0NGUtNGJiMi1iMzJjLTFmNWJmZjc3NTJmMiJdLCJhbWJlci1mYWl0aGZ1bC1zZXJ2aWNlLWlkcyI6WyI2YmFhNjMwMS0zMWVlLTQ1NmMtOWEzOC1lMjc5YWM3ZjZkNmEiLCJkZTU3NDU5ZC1mYjU2LTRhODgtYmU1ZC02ZjJhMzcwMzkzOWMiXSwiYW1iZXJfdGNiX3N0YXR1cyI6Ik9LIiwiYW1iZXJfZXZpZGVuY2VfdHlwZSI6IlNHWCIsImFtYmVyX3NpZ25lZF9ub25jZSI6dHJ1ZSwiYW1iZXJfY3VzdG9tX3BvbGljeSI6e30sInZlciI6IjEuMCIsImV4cCI6MTY2NDQ0NzA5MywianRpIjoiMzc5ZjBiMDctNTUzNy00YzdhLWFlNTAtODk2YTU1ZjIzNzY2IiwiaWF0IjoxNjY0NDQ2NzYzLCJpc3MiOiJBUyBBdHRlc3RhdGlvbiBUb2tlbiBJc3N1ZXIifQ.X2UDvraRVzAJpC1G1WAK2Qbx64d9WI5T_AKAq1lK5VAjEf409y5fZPxkBdZ-fGYt653nQ5Ah0-jkFRt0Yo7B2cxNmDWn61mMW9yYtt_55qHcbuDX5x4a-7MVawWjS1gLzY7qddmpzoIhwrx575c5JoQjG4qybDejRufUxhvu_XOOxSfhyh4JGRxBYNX19ZGeIbHtE3mfAXqg6qphZFfClIQLdlU-wGbefyN5mwpTK0T3eQ9Tlt0zZFrcv7lNIAPHHB52Ke9R7qdFEoVNNX-8YFMzk4gQyZdzYJS7Q3ElhQYXhBWnY5iwquftQztQcfydJL8o1OC-Ru0s-keF7OBaxABHcv5OhUKlVc44zaBnekP9lTzRCINYnVK67KxyAHsgXkK19UiX6v1FYxdcmdwZgNn5OkCwxiAMLgB8_CQku6q4aeyhCMo4acD1xKd6kkfgYQDxehLbGV6weT4E60omx6UFE13L9yANNNoWtzy0A4PJsiw0tbRPYYO8ehZ8Vgrb9sc00cdqG7_7ok-iivuxklaaSuzrY8VtkGw9T8g0w__fJ0X2KCMPcl3XfNidhOGxJ9402ff93X-QY3dHyaLOqmtJK0vlQ0vuoThseBSOezETalhFCuh-JUYZskokQ21fDPs2xDiytKubxqzrJVF1G1n1AVNWlIZXPXLyoXANS4s"
	tokenMalformedJKU = "ewogICJhbGciOiAiUFMzODQiLAogICJqa3UiOiAiYm9ndXNcbmJhc2VcblVSTCIsCiAgImtpZCI6ICIzYzI0MThiNWU2OWUyNjQ0YjkxNjcyZmY2MDU2NjY0ZDkyNDIzNGYwIiwKICAidHlwIjogIkpXVCIKfQ.eyJhbWJlcl90cnVzdF9zY29yZSI6MTAsImFtYmVyX3JlcG9ydF9kYXRhIjoiZWZmNWEyYTExNDg2N2FhOTQ0NjIwYzQ4Y2Q4NjcwNDZkYmY2ZjdmY2JmODQ5YTliNjZhNzg3MjJmNGRjZDdjOTAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLCJhbWJlcl90ZWVfaGVsZF9kYXRhIjoiQVFBQkFLMkxyUDhSZ1ROMW5naFJZWTBKQnZZQ1M3K2JCcVhoRjIzY2JkVFVRR3F3MEl2Wm9NYkpySUNmQXJ1MjVWWWpKbkZaS0Vvb1hRWmhPZUlBeTZNV0RpWEpmdDc0VnVUR25YZnNLUDk4bWNvZXBiQ2M4U1BJRFBsdkhTQy9QQWtlRzJUdlZ3QkhBbjcvcURVNXJwUENIS05xWDYweUd2SW95QjhrLzBJTzR5M2V0ekdvQjF5YVRPQ3Iyd1NCYmdUUkV1M3ppY3JJODFPL1RsK0FaWitVekFCdUxSUEgxdlBrelBQYVhCT21IN2Q4SnZXZ0RwSjhFenBNRitzakt4dXI1dkEraDNjamxDUG4yZjFjODhqZGIyNDgyUUJ2STZoanB4R2k0dWRUVVdJekdKdElKeElUbnNscThwTFpUekhnR3V2UEdYK2xkYWFKdUVPSGJBeFhDRW4zTnNGNHZvVjFSQ2I1OXBBMEI2NnZBd1RHZmFONE9pR205aGhiTk1NTnZNeGlhZmdGanJWWHpjc1BvUE5vN2hPd0dMcVJFdGUrMWkzZzlGNDBCK2hEZVV6elZhTU8zVkxHTUtEcDlUSDJqMytYSnRnU3p4dThOWlg1WEZVeGpSMlJINzV5d25vbnRNQStnaDZid1d1UUlWWWI2K0k3eHEzdWxOaUZldzZ4eWc9PSIsImFtYmVyX3NneF9tcmVuY2xhdmUiOiI4M2Y0ZTgxOTg2MWFkZWY2ZmZiMmE0ODY1ZWZlYTkzMzdiOTFlZDMwZmEzMzQ5MWIxN2YwZDVkOWU4MjA0NDEwIiwiYW1iZXJfc2d4X2lzX2RlYnVnZ2FibGUiOmZhbHNlLCJhbWJlcl9zZ3hfbXJzaWduZXIiOiI4M2Q3MTllNzdkZWFjYTE0NzBmNmJhZjYyYTRkNzc0MzAzYzg5OWRiNjkwMjBmOWM3MGVlMWRmYzA4YzdjZTllIiwiYW1iZXJfc2d4X2lzdnByb2RpZCI6MCwiYW1iZXJfc2d4X2lzdnN2biI6MCwiYW1iZXJfbWF0Y2hlZF9wb2xpY3lfaWRzIjpbImY0MzZjMzBhLTY0NGUtNGJiMi1iMzJjLTFmNWJmZjc3NTJmMiJdLCJhbWJlci1mYWl0aGZ1bC1zZXJ2aWNlLWlkcyI6WyI2YmFhNjMwMS0zMWVlLTQ1NmMtOWEzOC1lMjc5YWM3ZjZkNmEiLCJkZTU3NDU5ZC1mYjU2LTRhODgtYmU1ZC02ZjJhMzcwMzkzOWMiXSwiYW1iZXJfdGNiX3N0YXR1cyI6Ik9LIiwiYW1iZXJfZXZpZGVuY2VfdHlwZSI6IlNHWCIsImFtYmVyX3NpZ25lZF9ub25jZSI6dHJ1ZSwiYW1iZXJfY3VzdG9tX3BvbGljeSI6e30sInZlciI6IjEuMCIsImV4cCI6MTY2NDQ0NzA5MywianRpIjoiMzc5ZjBiMDctNTUzNy00YzdhLWFlNTAtODk2YTU1ZjIzNzY2IiwiaWF0IjoxNjY0NDQ2NzYzLCJpc3MiOiJBUyBBdHRlc3RhdGlvbiBUb2tlbiBJc3N1ZXIifQ.X2UDvraRVzAJpC1G1WAK2Qbx64d9WI5T_AKAq1lK5VAjEf409y5fZPxkBdZ-fGYt653nQ5Ah0-jkFRt0Yo7B2cxNmDWn61mMW9yYtt_55qHcbuDX5x4a-7MVawWjS1gLzY7qddmpzoIhwrx575c5JoQjG4qybDejRufUxhvu_XOOxSfhyh4JGRxBYNX19ZGeIbHtE3mfAXqg6qphZFfClIQLdlU-wGbefyN5mwpTK0T3eQ9Tlt0zZFrcv7lNIAPHHB52Ke9R7qdFEoVNNX-8YFMzk4gQyZdzYJS7Q3ElhQYXhBWnY5iwquftQztQcfydJL8o1OC-Ru0s-keF7OBaxABHcv5OhUKlVc44zaBnekP9lTzRCINYnVK67KxyAHsgXkK19UiX6v1FYxdcmdwZgNn5OkCwxiAMLgB8_CQku6q4aeyhCMo4acD1xKd6kkfgYQDxehLbGV6weT4E60omx6UFE13L9yANNNoWtzy0A4PJsiw0tbRPYYO8ehZ8Vgrb9sc00cdqG7_7ok-iivuxklaaSuzrY8VtkGw9T8g0w__fJ0X2KCMPcl3XfNidhOGxJ9402ff93X-QY3dHyaLOqmtJK0vlQ0vuoThseBSOezETalhFCuh-JUYZskokQ21fDPs2xDiytKubxqzrJVF1G1n1AVNWlIZXPXLyoXANS4s"
	//certs/crls no need to keep in hex form
	certHex             = "30820557308203bfa003020102020102300d06092a864886f70d01010d0500305b310b3009060355040613025553310b300906035504080c024341311a3018060355040a0c11496e74656c20436f72706f726174696f6e3123302106035504030c1a496e74656c20416d62657220415453205369676e696e67204341301e170d3233303131303033343230355a170d3233303730393033343230355a3060310b3009060355040613025553310b300906035504080c024341311a3018060355040a0c11496e74656c20436f72706f726174696f6e3128302606035504030c1f416d626572204174746573746174696f6e20546f6b656e205369676e696e6730820222300d06092a864886f70d01010105000382020f003082020a0282020100b833bb7f44d9a1521bf304c78b4080e3688d82c9fb5a8e9c513f314e7eeb4b87e92dfd04f029e96206a4a249e777c266fd69868dbc62d824261fa1c5656f811ea67b035155e806a75b858de937b65d0b9f2312d1bab91eb84697e064bee5fe63f5717c8aca3d50e075f1a8e8284f7cfee324a18060c9189bc1f92daf72cc8475151c1244e39513d8339aeb2cdcb39665065138356bdfc1c008e8ed382894fc662300b2caffcd52c06e739d7f5533414a7578b664156953d6b260dea206c8b59a02f60968d813bf75a8a0d8fb73f16d08bfcc88ec708da44ac7b6a227c81ff5e053439345f4a4e99fb0f2846630aab22123143486a705855113e81ba1aef52c31875a492f4dc1114be9254b2c86f2827c938add6ff35b5e112e7149132964abb2e4aabbf438cac65947ff38c171338be9323ecb1b101d0c2d6f38d6cb774de752c20f19569deb6f040943eb855225ca143de4265c8d199dec1f7f6d06b4dd6382ce22101f139533175972254bd782ad1f83e2ed3294611c6a9307cbe4d79e0db1b71e3bcee64f8af3ee5d0da52eb3ff9de4e0bc76e79c2cd5e58c5699bd8d755dad9ebe6f40a1c64e806c52ed7ee9bfe1e87f993d7d37cc7ef37d56aed50ba41dcecd52f35e83fa6b34d0ed6bc438c9a2e520f674977577856552b5f53ed9ba3083a92f81533d61d7321a3b355b8869c93a777661b635abf72395fb8f736ab7130203010001a381a030819d300c0603551d130101ff04023000301d0603551d0e04160414e4ddc5bd96e128a82465cb26e1bce1990d587898301f0603551d23041830168014438354f9cc18b8a30994672336bdbdf38c6372ab300b0603551d0f0404030204f030400603551d1f043930373035a033a031862f5552493a68747470733a2f2f616d6265722e696e74656c662e636f6d2f6174732d7369676e696e672d63612e63726c300d06092a864886f70d01010d0500038201810018f3d37036af59f676ff946549d4b9d17161c2194edf62123dd06935604f455fe50d0404138ad931d964fc41d2c7b6ba96dc6b3311e7ff274d77400e8eca05cb7c5bcaf3e48bce57733ce5d1730e6f035c02e69e139666976d74201c48bad05e57f15804c296a42c041b253d969e97d4a7bae032b6bec46eb0a30c0b57a660358663fd58f629410ddf20dee8450eecb5f4d3e2bb5a83787ab5b18c204b9cf9b76ff38c82e37ba4a7b7fc3cfeb52f5c74b67083c1151f512a4ed80795e6209d233f8fef198774aa1d2f33887219bde3e89d13823f96decc3a606876f1236e80b929364d7590359aa4846bda5674d96ec218db22e82001a514983334d7b3b4435f757728208757fde9f6a6988eedefa7602d8676ab23b0a4198431473068fc9dfc29935018fcfa0888992c37b805ccd0925dec8e8c3b062e9f8285b96a3de587a6f37dd790a26f8c575d54b1cd62fa50e780c36ffb600192e1c77b53cd7e7f6ac8a8bf18191c642c3ec3fe7d69fb980a3697549b9e33eb312c748a095293583059"
	invalidCRLCertHex   = "30820557308203bfa003020102020102300d06092a864886f70d01010d0500305b310b3009060355040613025553310b300906035504080c024341311a3018060355040a0c11496e74656c20436f72706f726174696f6e3123302106035504030c1a496e74656c20416d62657220415453205369676e696e67204341301e170d3233303431303137343732355a170d3233313030373137343732355a3060310b3009060355040613025553310b300906035504080c024341311a3018060355040a0c11496e74656c20436f72706f726174696f6e3128302606035504030c1f416d626572204174746573746174696f6e20546f6b656e205369676e696e6730820222300d06092a864886f70d01010105000382020f003082020a0282020100b833bb7f44d9a1521bf304c78b4080e3688d82c9fb5a8e9c513f314e7eeb4b87e92dfd04f029e96206a4a249e777c266fd69868dbc62d824261fa1c5656f811ea67b035155e806a75b858de937b65d0b9f2312d1bab91eb84697e064bee5fe63f5717c8aca3d50e075f1a8e8284f7cfee324a18060c9189bc1f92daf72cc8475151c1244e39513d8339aeb2cdcb39665065138356bdfc1c008e8ed382894fc662300b2caffcd52c06e739d7f5533414a7578b664156953d6b260dea206c8b59a02f60968d813bf75a8a0d8fb73f16d08bfcc88ec708da44ac7b6a227c81ff5e053439345f4a4e99fb0f2846630aab22123143486a705855113e81ba1aef52c31875a492f4dc1114be9254b2c86f2827c938add6ff35b5e112e7149132964abb2e4aabbf438cac65947ff38c171338be9323ecb1b101d0c2d6f38d6cb774de752c20f19569deb6f040943eb855225ca143de4265c8d199dec1f7f6d06b4dd6382ce22101f139533175972254bd782ad1f83e2ed3294611c6a9307cbe4d79e0db1b71e3bcee64f8af3ee5d0da52eb3ff9de4e0bc76e79c2cd5e58c5699bd8d755dad9ebe6f40a1c64e806c52ed7ee9bfe1e87f993d7d37cc7ef37d56aed50ba41dcecd52f35e83fa6b34d0ed6bc438c9a2e520f674977577856552b5f53ed9ba3083a92f81533d61d7321a3b355b8869c93a777661b635abf72395fb8f736ab7130203010001a381a030819d300c0603551d130101ff04023000301d0603551d0e04160414e4ddc5bd96e128a82465cb26e1bce1990d587898301f0603551d230418301680143bd92d080eb2c16881832d0d7436dc926e59c743300b0603551d0f0404030204f030400603551d1f043930373035a033a031862f5552493a68747470733a2f2f616d6265722e696e74656c662e636f6d2f6174732d7369676e696e672d63612e70656d300d06092a864886f70d01010d050003820181004279fb07bc3f336be663c838761dddb8ce5f7d28ae38800a707c4eb11c39434cb96fd741070ed1cb7becca26ab075deeb23667cfd533918fd3ccae394de3cc1ee07c800a0cb2ac7355beeb825a9a7fd791e63e8a058f260a7bf4bed9a1001154b3b9935fcc92f473aa86f933bc91e31677c4b5e08224e67e0b473ef8e82252a8430d6feec4cfd40ba0a64788346d29e1464f7e9cd497dcd35e6b561f8993d664f9d1cb08db555bd5510297cb766b7db0f032ab6ee57a495e485f87fc78fb5312948856b3c0c62b8056dc1308b40faf5bcc2fafdc871e464fb05398f4d47d35a9b39c0fdcb74ccea38bc41a821037a8e30190e998865fa8f0ae714aa4145e981a1601909eba2a9f5cd9584c7a4ae160396a692aa33d8fd686fe951e41d8ab14d01534c12477a5fae4ac715f4ee75da1e38f1689dfee2f4596da9ec069f000e80bb750a6aa2993bd049c1423ce1b677f5a7b6d28a962227507761acc156413d88baaf3bb3cb3c570aedbb11a7cb12f8329739abc44541d9d2811f0407816324433"
	invalidCACertHex    = "308204c53082032da003020102021436441d9104098f2aa434e0c9d36db2ce0cf80828300d06092a864886f70d01010d0500306a311c301a06035504030c13496e74656c20416d62657220526f6f74204341310b3009060355040613025553310b300906035504080c0243413114301206035504070c0b53616e746120436c617261311a3018060355040a0c11496e74656c20436f72706f726174696f6e301e170d3233303431303137343732315a170d3439313233303137343732315a306a311c301a06035504030c13496e74656c20416d62657220526f6f74204341310b3009060355040613025553310b300906035504080c0243413114301206035504070c0b53616e746120436c617261311a3018060355040a0c11496e74656c20436f72706f726174696f6e308201a2300d06092a864886f70d01010105000382018f003082018a0282018100caaef4384fe8e11855f33557e7d6a9bbba6e578eaba2f6f883e1582ba44b3da4a0980512a29e59b8c07a554488a11183aa94ee84a3540f4d9995431189c476e8e62ca83914c243916384a31a59c0a8b647931cfb7455927164d51071942754e286e9792b176396a18138d44b944fc73263a1d0064b723298802c81bf7a5fd1a4a773a328c687b3c050e7395929879ad5fd8b31e64e581b6015e38ab6aa511b3ac4b703e6a2d622fd4d21748694081294038cec2b0a2122addb78e8247e67ff4ce4c1ac93c1db2f24d2ac857cef6d431817d01360eec2249fb00e29a7195f72ab8b9ada0a8d09b571c666ab70c93a265a54e5fbba330cd53abf62091f42a7fe55da41c7ba0735865941fb65a4e2c714b3c6b1bb3822947e6263fbfd4c1f6e9d15d7d4b94cf6dcd0868a966a823c2dd41aee7c7c74c1fa3f12c7ead3032559f751065fa13664adf4db057f4810f52b41009c79e57535cbc268300e2addbd342efce224f2d87a47db89ba31fee519b873044e40abacfb3cb503cc75ac63b5e732550203010001a3633061300f0603551d130101ff040530030101ff301d0603551d0e04160414de4aa628f84dc5e960d5a850207fa684acfe7a56301f0603551d23041830168014de4aa628f84dc5e960d5a850207fa684acfe7a56300e0603551d0f0101ff040403020106300d06092a864886f70d01010d05000382018100774792455f8ff0e45605d3b4edbf744563a3cebaf93a29116d5b585b5e74134ade6524be5652d72e4cc57e535f9639747f9dd15f252d3423cf1f0dfa0516c8c219e6d7ee9a9859977cfc5f81e887da7f1e12f52a43bb093c01350e4c25dea4795a208bc883960d5c29cc5a016864c6e0ccb275257b43e3e0b549dd5c6181a50802413676d1887ba2b1df8dcf3d7918d46187fdf0b5a4c52b2613af55c3bf57af8c619e73801084a9ffaa302ecf0fd31fa423025b3fa9cf88debcaf798eb3ecf188dd7d3362d476de168e533689ec32e67d4dc94a10d1eb2492b020fa72fa1c700baad9be5917f81ae9e8775c1b8c9f19127048b93d5bfbf201ad19ce462e1f85bf81d7e79fe5521071e9abe6b11dea7540a2bc89791b0f09354068b5466d2bd2f0a23aad44b54e1114ad8ae61b9615f76bd120d431e0494bcbd6750f5d6444b79e6bd0163c4a83dd312ebd55cdd63d1aacd808ab5155012d22ea15ec40dcbc6a69643d2ffe4fd970abe003e9f54e26a7234c8df9a972a7104b14675b8e8e3129"
	crlHex              = "2d2d2d2d2d424547494e20583530392043524c2d2d2d2d2d0a4d4949434e5443426e6749424154414e42676b71686b69473977304241513046414442624d517377435159445651514745774a56557a454c4d416b47413155450a4341774351304578476a415942674e5642416f4d45556c756447567349454e76636e4276636d4630615739754d534d7749515944565151444442704a626e526c0a6243424262574a6c6369424256464d6755326c6e626d6c755a7942445152634e4d6a4d774e4445784d4459304d4455795768634e4d6a4d774e5445784d4459300a4d445579577141504d413077437759445652305542415143416841424d413047435371475349623344514542445155414134494267514148317351752f3466420a6334502b434e38592f62374861456c386b6f68323837666472735a2f5154522f446b79427648765835304751646e54394979724a6858475052767a69727652420a47385a63492f72556a326a4547644d71594d2f6d582b364c4c7a316c357551534367594b744a6e5351374647614e734c55744e414c35346d527455345246522b0a676b7566597134644638664c56646b366b6f6947567049356f452b597546356935627969373732464e5875703150733644474b785735776f38415259706d51610a5a79477a48524d546b4564543765567850785834755057546a7557646c76743153332b6d466d5861386b4a32776230434858685553342f33323770626f424f460a33512f56334c5951682b5154305671416a4f6c304e4b6b32626758387948587557692f4371642b7263734c2f4a7a4e4a6e733053373731742b636a79684368510a774761664a5655526c693350536c657853444c4b4a5957624f672b71616c4955352b456b582f4475324a5274417779496873726e456954566c7576394c39684c0a5565312b447467636849347556376b355961646f65416d6c5774743976597775496d313567682f727150716d72657141456378344c7a686c6d72586e35775a590a417745754e2b4664766152794e69644e59795967783244726645366778746c58616c6773564c694c75305838797056456c6436354538513d0a2d2d2d2d2d454e4420583530392043524c2d2d2d2d2d0a"
	validCertHex        = "30820557308203bfa003020102020102300d06092a864886f70d01010d0500305b310b3009060355040613025553310b300906035504080c024341311a3018060355040a0c11496e74656c20436f72706f726174696f6e3123302106035504030c1a496e74656c20416d62657220415453205369676e696e67204341301e170d3233303431313036343035325a170d3233313030383036343035325a3060310b3009060355040613025553310b300906035504080c024341311a3018060355040a0c11496e74656c20436f72706f726174696f6e3128302606035504030c1f416d626572204174746573746174696f6e20546f6b656e205369676e696e6730820222300d06092a864886f70d01010105000382020f003082020a0282020100b833bb7f44d9a1521bf304c78b4080e3688d82c9fb5a8e9c513f314e7eeb4b87e92dfd04f029e96206a4a249e777c266fd69868dbc62d824261fa1c5656f811ea67b035155e806a75b858de937b65d0b9f2312d1bab91eb84697e064bee5fe63f5717c8aca3d50e075f1a8e8284f7cfee324a18060c9189bc1f92daf72cc8475151c1244e39513d8339aeb2cdcb39665065138356bdfc1c008e8ed382894fc662300b2caffcd52c06e739d7f5533414a7578b664156953d6b260dea206c8b59a02f60968d813bf75a8a0d8fb73f16d08bfcc88ec708da44ac7b6a227c81ff5e053439345f4a4e99fb0f2846630aab22123143486a705855113e81ba1aef52c31875a492f4dc1114be9254b2c86f2827c938add6ff35b5e112e7149132964abb2e4aabbf438cac65947ff38c171338be9323ecb1b101d0c2d6f38d6cb774de752c20f19569deb6f040943eb855225ca143de4265c8d199dec1f7f6d06b4dd6382ce22101f139533175972254bd782ad1f83e2ed3294611c6a9307cbe4d79e0db1b71e3bcee64f8af3ee5d0da52eb3ff9de4e0bc76e79c2cd5e58c5699bd8d755dad9ebe6f40a1c64e806c52ed7ee9bfe1e87f993d7d37cc7ef37d56aed50ba41dcecd52f35e83fa6b34d0ed6bc438c9a2e520f674977577856552b5f53ed9ba3083a92f81533d61d7321a3b355b8869c93a777661b635abf72395fb8f736ab7130203010001a381a030819d300c0603551d130101ff04023000301d0603551d0e04160414e4ddc5bd96e128a82465cb26e1bce1990d587898301f0603551d23041830168014a749d7b9acab0d3d0674bcd1b6fdcf8af918eca5300b0603551d0f0404030204f030400603551d1f043930373035a033a031862f5552493a68747470733a2f2f616d6265722e696e74656c662e636f6d2f6174732d7369676e696e672d63612e63726c300d06092a864886f70d01010d0500038201810023dc69c857d1a91c98d9ccd752b4273ffe5eeba1836203c37bde0c25319f006c7d3c07e60c3a3010de2e45d9e620dfb5b84666fd35cf4405e0762b81b9bc6399f6cbde28539685016ca135169df59607f69544e87121ed2c4897f4815f8699b820dd8375cd8171d60502c74ff51293d0b1f7bb4c12bcc0c652319f63764a193cd5f03bde07411dff41cac76857349398656e7967d86e37d9d21dc84f7baff641754729b321ae18917841e12b944f246e42395856b4eea48038a819af986655a914578d382ee61c6bab6e17e57c24923f455a08e1842b3f8f3570901f85f1840ff68b732e890a59c0ee983a60167a47ada8d24563cb0bc35453b94a315f54bbf7d680e3740183b2050d41df707f98152eff4efa169dfb4c9716b96478e7c51bcfc4bc9ef524e0be49a274ae2f2440d1e6a38d9625c220a8994d6d69a4cfa4a9a48d98ff09504d69ecf7e4b561db016e511dce024fb820ff4193d791cd109a5bd6f1f3406b82d14a7d79a6809d454aaf129c9428cb708608025e850379ba633e6a"
	validInterCaCertHex = "308204d030820338a003020102020101300d06092a864886f70d01010d0500306a311c301a06035504030c13496e74656c20416d62657220526f6f74204341310b3009060355040613025553310b300906035504080c0243413114301206035504070c0b53616e746120436c617261311a3018060355040a0c11496e74656c20436f72706f726174696f6e301e170d3233303431313036343034355a170d3336313233303036343034355a305b310b3009060355040613025553310b300906035504080c024341311a3018060355040a0c11496e74656c20436f72706f726174696f6e3123302106035504030c1a496e74656c20416d62657220415453205369676e696e67204341308201a2300d06092a864886f70d01010105000382018f003082018a0282018100ac5a74938c1952036c8e8f920a751f9462eb423fe53984d06a8d3cd93c269d0dc7b741941b4c5e3068edbe3d4bc80bc7512bf07460d400956533e220b34ba72033f70f6d6e7bab3a27f35623834fed63560615e4240c3860039d8486438f50780a13ac31619416d681ebea5393cd90696b19a082bfe631f74a43ae1802e082dda0a7da12538ebb24c14bca3d407a222812c693a54ef99e58f02e37449d4b3d72b9d77c539374780dd4c6d5178c6b77ad29a788eaf80c0ca19018cd07c908779773540d6d1e34f14145bd408cdad0577dd6c9dfa9b3c5568e0ebd31ca9d55e251b6fdda44a56a73519b93fc766937dc35462abdf75a107a47487fe4f4f81e73d90505adb3097c811e80b60945022b56b568b01005e51fe8c9185e84dcca5233610f75ed978d1bb9ce73e8ed09cd011b8b1b43220d8a73e43f728db852566061131856126ddfc13024125d6b362f68039241ca88f34f4c525df334f6bd86a2e79a6109980cb7382c12787d7e6f6e28218ad813754bcfdece312f633f73eba579d90203010001a3818f30818c301d0603551d0e04160414a749d7b9acab0d3d0674bcd1b6fdcf8af918eca5301f0603551d23041830168014f06a5641e8cfdd949b7a57f26fa6d914e564878d300f0603551d130101ff040530030101ff30390603551d1f04323030302ea02ca02a86285552493a68747470733a2f2f616d6265722e696e74656c662e636f6d2f726f6f742d63612e63726c300d06092a864886f70d01010d050003820181007905d0af19abbec71149e2fc7fcb40b8826f287d2345d816ec9831e2db3c8cb23b47975d1d70d13f59fd5c53b6e664976379c564b1df0b9d8314427550d77ebb7430425257d11802782ab8fc65b04592b9bbbff8a6c5ded69c05a587370a225a008f4b7835586d76e2d156a1ffff44814596d869f5520524b7abeeb93318e63ffa698990abfae8a04fef62c4bea3f846deb66520a424ccc14fb24a82d6ac28623d28b117dec8274dfaaf31c84015b5bfac3047c97cd6b439c976a7333e5bdacd5ce7fca8ee3e7f9547c5c28136d639a9cb47b50081f7efcbf77ff32720a1cb2738949e31cd9ef5df082b84d4c143e5bb4db7a8633de7602fc6bcad40cc1060972336c08c7389bbbd51064fc78f8bad5317b0c954a20d6950e6f0e50886220110710dbb670b0c40d0ddd139ead1e4edd52159fb84c22b8f24033a0bb007f7c6cbd690e4c4eee7b0e6d435f377dcd2e3ca190a8e0dda44fb865bb1f84405f720be52b5afc6d310c7b46d68d721efd4bb1ae3e2f448a34e350bc7d5c7d677bde3f6"
)

func TestGetToken(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/appraisal/v1/attest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"token":"` + token + `"}`))
	})

	nonce := &Nonce{}
	evidence := &Evidence{}
	_, err := client.GetToken(nonce, nil, evidence)
	if err != nil {
		t.Errorf("GetToken returned unexpected error: %v", err)
	}
}

func TestGetToken_invalidToken(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/appraisal/v1/attest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid token`))
	})

	nonce := &Nonce{}
	evidence := &Evidence{}
	_, err := client.GetToken(nonce, nil, evidence)
	if err == nil {
		t.Errorf("GetToken returned nil, expected error")
	}
}

func TestVerifyToken_emptyToken(t *testing.T) {
	cfg := Config{
		Url: "https://custom-url/api/v1",
	}

	client, _ := New(&cfg)
	_, err := client.VerifyToken("")
	if err == nil {
		t.Error("VerifyToken returned nil, expected error")
	}
}

func TestVerifyToken_missingKID(t *testing.T) {
	cfg := Config{
		Url: "https://custom-url/api/v1",
	}

	client, _ := New(&cfg)
	_, err := client.VerifyToken(tokenMissingKID)
	if err == nil {
		t.Error("VerifyToken returned nil, expected error")
	}
}

func TestVerifyToken_invalidKID(t *testing.T) {
	cfg := Config{
		Url: "https://custom-url/api/v1",
	}

	client, _ := New(&cfg)
	_, err := client.VerifyToken(tokenInvalidKID)
	if err == nil {
		t.Error("VerifyToken returned nil, expected error")
	}
}

func TestVerifyToken_missingJKU(t *testing.T) {
	cfg := Config{
		Url: "https://custom-url/api/v1",
	}

	client, _ := New(&cfg)
	_, err := client.VerifyToken(tokenMissingJKU)
	if err == nil {
		t.Error("VerifyToken returned nil, expected error")
	}
}

func TestVerifyToken_invalidJKU(t *testing.T) {
	cfg := Config{
		Url: "https://custom-url/api/v1",
	}

	client, _ := New(&cfg)
	_, err := client.VerifyToken(tokenInvalidJKU)
	if err == nil {
		t.Error("VerifyToken returned nil, expected error")
	}
}

func TestVerifyToken_malformedJKU(t *testing.T) {
	cfg := Config{
		Url: "https://custom-url/api/v1",
	}

	client, _ := New(&cfg)
	_, err := client.VerifyToken(tokenMalformedJKU)
	if err == nil {
		t.Error("VerifyToken returned nil, expected error")
	}
}
func TestGetCRLObject_emptyCRLURL(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	var emptyCRLArry []string
	_, err := client.GetCRL(emptyCRLArry)
	if err == nil {
		t.Errorf("verifyCRL returned nil, expected error")
	}
}
func TestGetCRLObject_validCRLUrl(t *testing.T) {
	client, mux, serverUrl, teardown := setup()
	defer teardown()

	crlBytes, _ := hex.DecodeString(crlHex)
	crlUrl := serverUrl + "/ats.crl"
	mux.HandleFunc("/ats.crl", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(crlBytes)
	})

	_, err := client.GetCRL([]string{crlUrl})
	if err != nil {
		t.Errorf("GetCRL returned err,  expected nil")
	}
}

func TestGetCRLObject_invalidCRL(t *testing.T) {

	client, mux, serverUrl, teardown := setup()
	defer teardown()

	crlUrl := serverUrl + "/ats.crl"
	mux.HandleFunc("/ats.crl", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalidcrl"))
	})

	_, err := client.GetCRL([]string{crlUrl})
	if err == nil {
		t.Errorf("GetCRL returned nil,  expected error")
	}
}

func TestVerifyCRL_nullCerts(t *testing.T) {
	_, _, _, teardown := setup()
	defer teardown()

	var leafCert *x509.Certificate
	var interCaCert *x509.Certificate

	err := verifyCRL(nil, leafCert, interCaCert)
	if err == nil {
		t.Errorf("verifyCRL returned nil, expected error")
	}
}

func TestVerifyCRL_nilArguments(t *testing.T) {
	_, _, _, teardown := setup()
	defer teardown()

	certBytes, _ := hex.DecodeString(certHex)
	leafCert, _ := x509.ParseCertificate(certBytes)
	interCaCert, _ := x509.ParseCertificate(certBytes)

	err := verifyCRL(nil, leafCert, interCaCert)
	if err == nil {
		t.Errorf("verifyCRL returned nil, expected error")
	}
}
func TestVerifyCRL_validCertAndCrl(t *testing.T) {
	_, _, _, teardown := setup()
	defer teardown()
	crlBytes, _ := hex.DecodeString(crlHex)
	certBytes, _ := hex.DecodeString(validCertHex)
	caCertBytes, _ := hex.DecodeString(validInterCaCertHex)

	leafCert, _ := x509.ParseCertificate(certBytes)
	interCaCert, _ := x509.ParseCertificate(caCertBytes)

	block, _ := pem.Decode([]byte(crlBytes))
	crl, _ := x509.ParseRevocationList(block.Bytes)

	crl.NextUpdate = time.Now().AddDate(0, 0, 3)

	err := verifyCRL(crl, leafCert, interCaCert)
	if err != nil {
		t.Errorf("verifyCRL returned nil, expected error")
	}
}

func TestVerifyCRL_invalidCACert(t *testing.T) {

	_, _, _, teardown := setup()
	defer teardown()
	crlBytes, _ := hex.DecodeString(crlHex)
	certBytes, _ := hex.DecodeString(invalidCRLCertHex)
	caCertBytes, _ := hex.DecodeString(invalidCACertHex)
	leafCert, _ := x509.ParseCertificate(certBytes)
	interCaCert, _ := x509.ParseCertificate(caCertBytes)

	block, _ := pem.Decode([]byte(crlBytes))
	crl, _ := x509.ParseRevocationList(block.Bytes)

	err := verifyCRL(crl, leafCert, interCaCert)
	if err == nil {
		t.Errorf("verifyCRL returned nil, expected error")
	}
}
