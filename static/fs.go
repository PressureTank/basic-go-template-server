// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package static

import "embed"

//go:embed *

// FS defines a static filesystem that can be served by the server.
var FS embed.FS
