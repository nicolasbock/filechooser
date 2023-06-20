package main

import (
	"testing"
)

func compareFileList(a, b Files) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestReadFiles1(t *testing.T) {
	want := Files{
		{
			name:   "001-b1f16d37-71ee-4b69-928a-695da207b1ab.txt",
			path:   "artifacts/b/b/001-b1f16d37-71ee-4b69-928a-695da207b1ab.txt",
			md5sum: "e28695372b915ded0520febd37f4f1fc",
		},
		{
			name:   "002-769f1d04-5337-4728-8878-bd96318573db.txt",
			path:   "artifacts/b/b/002-769f1d04-5337-4728-8878-bd96318573db.txt",
			md5sum: "8e5920fa2f3f4f2eda08baa479ee8de9",
		},
		{
			name:   "003-97cb0609-0bfb-44a0-bfaf-9680001ac444.txt",
			path:   "artifacts/b/b/003-97cb0609-0bfb-44a0-bfaf-9680001ac444.txt",
			md5sum: "905b43b0faa458934e1b82fbdc56c025",
		},
		{
			name:   "004-10b91a28-7ac5-4f06-a986-af8f1d006a52.txt",
			path:   "artifacts/b/b/004-10b91a28-7ac5-4f06-a986-af8f1d006a52.txt",
			md5sum: "9a8c188e031bf809ee34389d6acd3f40",
		},
		{
			name:   "005-5e5c52ac-b24e-4e26-8078-aa6a90a03173.txt",
			path:   "artifacts/b/b/005-5e5c52ac-b24e-4e26-8078-aa6a90a03173.txt",
			md5sum: "04b69006cbe3d2e9490619bca0aa8448",
		},
		{
			name:   "006-e7d67845-31f6-464f-9c41-775383f4d8af.txt",
			path:   "artifacts/b/b/006-e7d67845-31f6-464f-9c41-775383f4d8af.txt",
			md5sum: "d014e983fd26eb7dd54f4b6e3a1d1673",
		},
		{
			name:   "007-d9714bbd-ad5a-4444-804f-c41003791899.txt",
			path:   "artifacts/b/b/007-d9714bbd-ad5a-4444-804f-c41003791899.txt",
			md5sum: "eb9dbcec3909c274147ec3124c66efc7",
		},
		{
			name:   "008-7cef8d7c-a80a-4aa7-b18d-f0c9c4c16514.txt",
			path:   "artifacts/b/b/008-7cef8d7c-a80a-4aa7-b18d-f0c9c4c16514.txt",
			md5sum: "f33e072d59b41657f189c313bfb919f0",
		},
		{
			name:   "009-731462cf-397e-40d7-9f06-377eca29aa0f.txt",
			path:   "artifacts/b/b/009-731462cf-397e-40d7-9f06-377eca29aa0f.txt",
			md5sum: "62afbb676f639a45ae70f6e59242caac",
		},
		{
			name:   "010-b42f3393-bfbc-4d74-8806-5f4b19683905.txt",
			path:   "artifacts/b/b/010-b42f3393-bfbc-4d74-8806-5f4b19683905.txt",
			md5sum: "a3d555055a0ce711a7ebcf21249e1972",
		},
		{
			name:   "011-5088df35-4118-4234-bf51-2eef040aac00.txt",
			path:   "artifacts/b/b/011-5088df35-4118-4234-bf51-2eef040aac00.txt",
			md5sum: "06657110877c0ebf825cab0fbb2e61f2",
		},
		{
			name:   "012-d515db3c-7e80-4496-ae9f-8cef6fbf2fb9.txt",
			path:   "artifacts/b/b/012-d515db3c-7e80-4496-ae9f-8cef6fbf2fb9.txt",
			md5sum: "01e79909acd34213829166c2f41e0095",
		},
		{
			name:   "013-700ef00d-c4cc-49e7-9ce1-7ba1a148eb71.txt",
			path:   "artifacts/b/b/013-700ef00d-c4cc-49e7-9ce1-7ba1a148eb71.txt",
			md5sum: "7fd1a416835e12c5bb0aebf617079e75",
		},
		{
			name:   "014-4dd8352f-349f-4c85-8937-4453ce502572.txt",
			path:   "artifacts/b/b/014-4dd8352f-349f-4c85-8937-4453ce502572.txt",
			md5sum: "9ff90815d210f8efbe7c31df14380de3",
		},
		{
			name:   "015-101ae6df-8c92-4384-b187-98f9fbef8a76.txt",
			path:   "artifacts/b/b/015-101ae6df-8c92-4384-b187-98f9fbef8a76.txt",
			md5sum: "101ecd9206be449d8ea23960218deb43",
		},
		{
			name:   "016-bc7dd8c3-1d45-466f-93e8-07220831ce30.txt",
			path:   "artifacts/b/b/016-bc7dd8c3-1d45-466f-93e8-07220831ce30.txt",
			md5sum: "477ccfe5caf23b508a00cec22e752f32",
		},
		{
			name:   "017-521c7a3b-a1da-4b29-bebd-bf8be1f7667c.txt",
			path:   "artifacts/b/b/017-521c7a3b-a1da-4b29-bebd-bf8be1f7667c.txt",
			md5sum: "55b95c0c59fe29f44949fab27be87bb1",
		},
		{
			name:   "018-0506ca48-7a60-4177-a67b-295bace5bfcd.txt",
			path:   "artifacts/b/b/018-0506ca48-7a60-4177-a67b-295bace5bfcd.txt",
			md5sum: "82fbbd4fab5527405cb84241e871e750",
		},
		{
			name:   "019-74f8eb80-95be-4e92-b3e3-8726d4abbf01.txt",
			path:   "artifacts/b/b/019-74f8eb80-95be-4e92-b3e3-8726d4abbf01.txt",
			md5sum: "15be5154d498296467df4ae5f7d6e822",
		},
		{
			name:   "020-aa150c61-238f-439d-bf0e-aa21636ef580.txt",
			path:   "artifacts/b/b/020-aa150c61-238f-439d-bf0e-aa21636ef580.txt",
			md5sum: "f7b41e826d47ef05e4f8e6667615302d",
		},
	}
	got := readFiles([]string{"artifacts/b/b"})
	if !compareFileList(want, got) {
		t.Errorf("want: %s != got: %s\n", want, got)
	}
}

func TestReadFiles2(t *testing.T) {
	want := Files{
		{
			name:   "001-44ff8d85-40bd-44c4-9625-cf676d99497a.txt",
			path:   "artifacts/a/a/001-44ff8d85-40bd-44c4-9625-cf676d99497a.txt",
			md5sum: "2162f63c07418fb6230c097167c8343e",
		},
		{
			name:   "002-f269e109-852b-4c82-882b-51f2aae97591.txt",
			path:   "artifacts/a/a/002-f269e109-852b-4c82-882b-51f2aae97591.txt",
			md5sum: "d1c8f7cb6f6434e633c75775e9fc5cd1",
		},
		{
			name:   "003-c665a6d7-5f07-4899-91af-2c2c22b1d73a.txt",
			path:   "artifacts/a/a/003-c665a6d7-5f07-4899-91af-2c2c22b1d73a.txt",
			md5sum: "766484e4e8f9063713f973b2fca0488e",
		},
		{
			name:   "004-23601d7f-43d4-44b7-8f41-15d1d2a99f4e.txt",
			path:   "artifacts/a/a/004-23601d7f-43d4-44b7-8f41-15d1d2a99f4e.txt",
			md5sum: "8928c730598c7068523160139a23ab8b",
		},
		{
			name:   "005-8dc11b52-fbba-4cf9-a27a-828bbc6d240a.txt",
			path:   "artifacts/a/a/005-8dc11b52-fbba-4cf9-a27a-828bbc6d240a.txt",
			md5sum: "e7ea83e581c1f0ae0e8636496fb5c272",
		},
		{
			name:   "006-32463186-bc6e-463d-aaaf-bc16318f91fd.txt",
			path:   "artifacts/a/a/006-32463186-bc6e-463d-aaaf-bc16318f91fd.txt",
			md5sum: "0bb00d95ea1f4d007dbcb6827191a208",
		},
		{
			name:   "007-1b041b6d-fd5a-4a18-a013-587ed0ddca69.txt",
			path:   "artifacts/a/a/007-1b041b6d-fd5a-4a18-a013-587ed0ddca69.txt",
			md5sum: "348a5eb81d7fba7ccb44b65b63ae3f3d",
		},
		{
			name:   "008-b2997063-2f08-4c04-ad2a-8ec0a126412a.txt",
			path:   "artifacts/a/a/008-b2997063-2f08-4c04-ad2a-8ec0a126412a.txt",
			md5sum: "3a3cea097921e6765b4abdcd21a7ff37",
		},
		{
			name:   "009-d10451ad-1d76-4fc1-ab4d-625514aecf39.txt",
			path:   "artifacts/a/a/009-d10451ad-1d76-4fc1-ab4d-625514aecf39.txt",
			md5sum: "d1b9508f4e177c0749fcbaa8b0fd8ae7",
		},
		{
			name:   "010-a889aec9-3576-4ebb-9895-5899843d12b3.txt",
			path:   "artifacts/a/a/010-a889aec9-3576-4ebb-9895-5899843d12b3.txt",
			md5sum: "b4193072cac90dde02c6d3ec6725a713",
		},
		{
			name:   "011-865945f6-b52e-42f0-b3a5-d8e6e279f15f.txt",
			path:   "artifacts/a/a/011-865945f6-b52e-42f0-b3a5-d8e6e279f15f.txt",
			md5sum: "8c5c86aa0dd3389161843d71d64b639a",
		},
		{
			name:   "012-d97889ed-a578-4cad-a6c0-fafdde8ad472.txt",
			path:   "artifacts/a/a/012-d97889ed-a578-4cad-a6c0-fafdde8ad472.txt",
			md5sum: "57dfd7f029da2be8b7ca9a702828f7c7",
		},
		{
			name:   "013-58c532fa-7ca2-4275-8223-44ea81ac88a9.txt",
			path:   "artifacts/a/a/013-58c532fa-7ca2-4275-8223-44ea81ac88a9.txt",
			md5sum: "5a99027812ca0163121497419e1de560",
		},
		{
			name:   "014-5f10d285-e232-4d6a-98af-11dda4954bf9.txt",
			path:   "artifacts/a/a/014-5f10d285-e232-4d6a-98af-11dda4954bf9.txt",
			md5sum: "007d546792608bf6f00ac4e0b91d9252",
		},
		{
			name:   "015-80e193d0-92d3-4ff0-b92f-2b71a96a0322.txt",
			path:   "artifacts/a/a/015-80e193d0-92d3-4ff0-b92f-2b71a96a0322.txt",
			md5sum: "2c2378b77233fd506c34d5429e195ad8",
		},
		{
			name:   "016-69e14c25-9610-412f-b704-a569f9f2077e.txt",
			path:   "artifacts/a/a/016-69e14c25-9610-412f-b704-a569f9f2077e.txt",
			md5sum: "591c66fe740e27871d799b2bf9e657df",
		},
		{
			name:   "017-36eb0436-3edd-4d1e-853d-9b94b9a147c6.txt",
			path:   "artifacts/a/a/017-36eb0436-3edd-4d1e-853d-9b94b9a147c6.txt",
			md5sum: "be0463db7c8d7be8caa0f98799c9e60c",
		},
		{
			name:   "018-cc40ff4b-9bb2-4990-b806-fb9fc20deab4.txt",
			path:   "artifacts/a/a/018-cc40ff4b-9bb2-4990-b806-fb9fc20deab4.txt",
			md5sum: "a594ebc131193a99c8fb422d06adcbc0",
		},
		{
			name:   "019-f3291097-bf25-41c6-9140-a41befbd3a88.txt",
			path:   "artifacts/a/a/019-f3291097-bf25-41c6-9140-a41befbd3a88.txt",
			md5sum: "007053297b251eda09bd60057034fb67",
		},
		{
			name:   "020-f3620c58-e14b-461a-8ea8-6948fbcdf8f8.txt",
			path:   "artifacts/a/a/020-f3620c58-e14b-461a-8ea8-6948fbcdf8f8.txt",
			md5sum: "7ff76395ad9fb0823f0ea1aecfec4d92",
		},
	}
	got := readFiles([]string{"artifacts/a/a"})
	if !compareFileList(want, got) {
		t.Errorf("got: %s, want: %s\n", got, want)
	}
}
