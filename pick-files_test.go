package main

import (
	"testing"
	"time"
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
			name:   "001-0ef0949e-9dae-44b3-844c-225950f174a1.jpg",
			path:   "artifacts/b/b/001-0ef0949e-9dae-44b3-844c-225950f174a1.jpg",
			md5sum: "68656be304743ea0ca18925844b2b633",
		},
		{
			name:   "001-1b73be01-b0b2-43e7-bacd-d4e14fab9c59.jpg",
			path:   "artifacts/b/b/001-1b73be01-b0b2-43e7-bacd-d4e14fab9c59.jpg",
			md5sum: "f9ec6b030b7f48cc167c6764c4b06f37",
		},
		{
			name:   "001-4c96bdcb-677e-43fc-9f78-b19ce0412dc3.png",
			path:   "artifacts/b/b/001-4c96bdcb-677e-43fc-9f78-b19ce0412dc3.png",
			md5sum: "faa8b000ff8d2e3234a6c00e7a0f6bcf",
		},
		{
			name:   "001-654cf09e-f0ec-420b-937f-f6ecf3f3cb2e.tif",
			path:   "artifacts/b/b/001-654cf09e-f0ec-420b-937f-f6ecf3f3cb2e.tif",
			md5sum: "5565a780891f4a4f0eb52887b4ce8214",
		},
		{
			name:   "001-8a34159b-7721-4842-a4ac-9ac43ea955f3.png",
			path:   "artifacts/b/b/001-8a34159b-7721-4842-a4ac-9ac43ea955f3.png",
			md5sum: "3888cac6d13d0346f787dc28c9e2350a",
		},
		{
			name:   "001-9b3e4593-c8b2-4cb0-977b-9d1ee4c70d2f.tif",
			path:   "artifacts/b/b/001-9b3e4593-c8b2-4cb0-977b-9d1ee4c70d2f.tif",
			md5sum: "62b2538bf5170a040637a1ebb8419224",
		},
		{
			name:   "002-035f1f87-f120-4efa-9806-2227658133b7.jpg",
			path:   "artifacts/b/b/002-035f1f87-f120-4efa-9806-2227658133b7.jpg",
			md5sum: "00acb4c0ccd76d6fdbac615506ee129e",
		},
		{
			name:   "002-0cbab6e4-190f-475c-a0c5-3369cfe35fda.jpg",
			path:   "artifacts/b/b/002-0cbab6e4-190f-475c-a0c5-3369cfe35fda.jpg",
			md5sum: "b77fe064b8d65cc27fa5850ba43a5879",
		},
		{
			name:   "002-78a80bea-5b56-4cb9-a09a-36ffbc4aaec6.tif",
			path:   "artifacts/b/b/002-78a80bea-5b56-4cb9-a09a-36ffbc4aaec6.tif",
			md5sum: "cb0efea2b9441e5a07bfa5fe825dbd1f",
		},
		{
			name:   "002-a1ba4c72-f8aa-4f28-a0a8-52581ab4afb6.png",
			path:   "artifacts/b/b/002-a1ba4c72-f8aa-4f28-a0a8-52581ab4afb6.png",
			md5sum: "c0f881d1633e0f1b94cc1902b588badf",
		},
		{
			name:   "002-c8d78dc0-e6d8-4e3b-b009-f90f9338078d.tif",
			path:   "artifacts/b/b/002-c8d78dc0-e6d8-4e3b-b009-f90f9338078d.tif",
			md5sum: "0b7b494d5e562eda72f2f4fada3a058d",
		},
		{
			name:   "002-d4e79a91-acc6-4ede-a9c3-9c7a9f2fa011.png",
			path:   "artifacts/b/b/002-d4e79a91-acc6-4ede-a9c3-9c7a9f2fa011.png",
			md5sum: "bd49757c9f3a26558eab170b4b636b0e",
		},
		{
			name:   "003-2b54f9e6-1094-4431-822d-57fb59ee644b.jpg",
			path:   "artifacts/b/b/003-2b54f9e6-1094-4431-822d-57fb59ee644b.jpg",
			md5sum: "e96c73aecbae596c433aefd9f6c369a4",
		},
		{
			name:   "003-412757c5-17d7-4b2d-b82e-2fbc2f3e5b73.tif",
			path:   "artifacts/b/b/003-412757c5-17d7-4b2d-b82e-2fbc2f3e5b73.tif",
			md5sum: "7fe8b6fceebd361240597f22ed5bc4aa",
		},
		{
			name:   "003-5edabc9c-969e-4379-a32d-1cfc1fa48398.jpg",
			path:   "artifacts/b/b/003-5edabc9c-969e-4379-a32d-1cfc1fa48398.jpg",
			md5sum: "758085dcd279ef9e73918f9738517894",
		},
		{
			name:   "003-6848877d-5968-4401-ad59-6c0a525042b3.png",
			path:   "artifacts/b/b/003-6848877d-5968-4401-ad59-6c0a525042b3.png",
			md5sum: "3a6d31a52b76f17513ff3e2cd630af92",
		},
		{
			name:   "003-8e839208-344a-4c25-8d05-9a4f1ea7081c.png",
			path:   "artifacts/b/b/003-8e839208-344a-4c25-8d05-9a4f1ea7081c.png",
			md5sum: "3d3b844160f505874955edd549b8f9c3",
		},
		{
			name:   "003-9cb3b0b7-0d52-4b0f-a3a6-f44c10392da5.tif",
			path:   "artifacts/b/b/003-9cb3b0b7-0d52-4b0f-a3a6-f44c10392da5.tif",
			md5sum: "d96eaa1cec5606a08ff483de133f3c59",
		},
		{
			name:   "004-33d02223-b6d8-46ec-9f08-363d38f38e60.tif",
			path:   "artifacts/b/b/004-33d02223-b6d8-46ec-9f08-363d38f38e60.tif",
			md5sum: "8c9feaa91683e0d016a39ac4ec52a4fc",
		},
		{
			name:   "004-b2721ec5-14be-4ffc-9b3e-747efee7fed9.png",
			path:   "artifacts/b/b/004-b2721ec5-14be-4ffc-9b3e-747efee7fed9.png",
			md5sum: "8dba93db8e8f7ce771ab150c3fa55a8b",
		},
		{
			name:   "004-b323a874-c293-418a-a520-bb57a505c992.jpg",
			path:   "artifacts/b/b/004-b323a874-c293-418a-a520-bb57a505c992.jpg",
			md5sum: "cf69416c5ccea3c47b594eb4df07f811",
		},
		{
			name:   "004-d2dccc6a-c245-4eaa-b63e-2ed679a975db.tif",
			path:   "artifacts/b/b/004-d2dccc6a-c245-4eaa-b63e-2ed679a975db.tif",
			md5sum: "4b7409aa9484e6535e0380fe700eabc6",
		},
		{
			name:   "004-d9722e43-6ab6-4de2-a691-d154dba39df6.png",
			path:   "artifacts/b/b/004-d9722e43-6ab6-4de2-a691-d154dba39df6.png",
			md5sum: "2e3e74f33c2d4e92f643ccda49d485ad",
		},
		{
			name:   "004-fa3c6b66-8268-4622-8102-62e58bcda47c.jpg",
			path:   "artifacts/b/b/004-fa3c6b66-8268-4622-8102-62e58bcda47c.jpg",
			md5sum: "c0e961a675be00b31a4e834be3a074e3",
		},
		{
			name:   "005-0c553984-77d2-41cc-9d06-b1dfb4f2cfc9.tif",
			path:   "artifacts/b/b/005-0c553984-77d2-41cc-9d06-b1dfb4f2cfc9.tif",
			md5sum: "b6f857ff3ce86bd5cd376e0cb6c6413c",
		},
		{
			name:   "005-4e9fc8d4-98a3-41b3-a679-b63df350e665.tif",
			path:   "artifacts/b/b/005-4e9fc8d4-98a3-41b3-a679-b63df350e665.tif",
			md5sum: "7009f449a00f3b461c589909208b84ac",
		},
		{
			name:   "005-6e8fd924-ba46-45fb-b33e-e7139032482b.png",
			path:   "artifacts/b/b/005-6e8fd924-ba46-45fb-b33e-e7139032482b.png",
			md5sum: "51dc44d1e9451922ff2053b399d44629",
		},
		{
			name:   "005-a8de4486-0733-42b3-a035-9f069159068e.png",
			path:   "artifacts/b/b/005-a8de4486-0733-42b3-a035-9f069159068e.png",
			md5sum: "e35bd9ca2bd12d7eb5eaaf8048169ec9",
		},
		{
			name:   "005-f6d94b28-69aa-436b-a8e0-43f238615ce2.jpg",
			path:   "artifacts/b/b/005-f6d94b28-69aa-436b-a8e0-43f238615ce2.jpg",
			md5sum: "84e828fed348e3b01f7ed56d87f2d999",
		},
		{
			name:   "005-f6e91086-940b-4e86-8986-f7cf9d9bc0c0.jpg",
			path:   "artifacts/b/b/005-f6e91086-940b-4e86-8986-f7cf9d9bc0c0.jpg",
			md5sum: "77efe5b443fbf39aa83359b875d2ec3e",
		},
		{
			name:   "006-244bb2c2-c0e5-4f2f-b7b2-bdda25f91d67.jpg",
			path:   "artifacts/b/b/006-244bb2c2-c0e5-4f2f-b7b2-bdda25f91d67.jpg",
			md5sum: "f0a820b3db8bef4971d1ef67e2bf276b",
		},
		{
			name:   "006-5b352373-feb0-490b-a4c8-84edbe7a1b68.png",
			path:   "artifacts/b/b/006-5b352373-feb0-490b-a4c8-84edbe7a1b68.png",
			md5sum: "9e2482dccd0cb5e5e29ca8c20aefe381",
		},
		{
			name:   "006-6577b65d-cd1d-4855-ba21-564f2c5529a7.jpg",
			path:   "artifacts/b/b/006-6577b65d-cd1d-4855-ba21-564f2c5529a7.jpg",
			md5sum: "0155fe31436ad0621942817d3f1ddd24",
		},
		{
			name:   "006-6e9aa5e7-7334-40bd-b130-cb22bb6cdf14.tif",
			path:   "artifacts/b/b/006-6e9aa5e7-7334-40bd-b130-cb22bb6cdf14.tif",
			md5sum: "415d933fe4fb6863bd81b972127b297e",
		},
		{
			name:   "006-d7a91e88-45a1-4fdc-9da4-ef39ef31a60e.tif",
			path:   "artifacts/b/b/006-d7a91e88-45a1-4fdc-9da4-ef39ef31a60e.tif",
			md5sum: "dae3a461d2058d16f01b1fc97895fc62",
		},
		{
			name:   "006-fd4c6143-4bc9-4630-9a31-785d51db7941.png",
			path:   "artifacts/b/b/006-fd4c6143-4bc9-4630-9a31-785d51db7941.png",
			md5sum: "b066f11d139ca3edbf461d32393176b7",
		},
		{
			name:   "007-556fb447-7d0d-4951-8e9d-6a6034183d25.jpg",
			path:   "artifacts/b/b/007-556fb447-7d0d-4951-8e9d-6a6034183d25.jpg",
			md5sum: "745baf700f431e9aec87363b23e10688",
		},
		{
			name:   "007-ae1b04cd-dca0-48fe-b9d7-ef6b1495caf4.jpg",
			path:   "artifacts/b/b/007-ae1b04cd-dca0-48fe-b9d7-ef6b1495caf4.jpg",
			md5sum: "a5186004467f2ea38d4153577f345ccd",
		},
		{
			name:   "007-d08bd0fc-56ee-4274-b06e-2569eea26a9b.png",
			path:   "artifacts/b/b/007-d08bd0fc-56ee-4274-b06e-2569eea26a9b.png",
			md5sum: "cf6a855171e73119c91eaab02d81b228",
		},
		{
			name:   "007-ed516b60-ba1f-481d-829e-d3186b1b53df.tif",
			path:   "artifacts/b/b/007-ed516b60-ba1f-481d-829e-d3186b1b53df.tif",
			md5sum: "2fff48e15da324fa6d8d5e0821288822",
		},
		{
			name:   "007-efad9dd5-a00b-4858-8188-aa75b9575193.tif",
			path:   "artifacts/b/b/007-efad9dd5-a00b-4858-8188-aa75b9575193.tif",
			md5sum: "9d111979a91f9cf01903e33733a23252",
		},
		{
			name:   "007-f3d63573-3643-44e9-a438-91076efa5264.png",
			path:   "artifacts/b/b/007-f3d63573-3643-44e9-a438-91076efa5264.png",
			md5sum: "9d32e916bb6815dc454d64a294ed6b7a",
		},
		{
			name:   "008-06b46cb6-9b5a-4132-acd3-55a4fadcb248.jpg",
			path:   "artifacts/b/b/008-06b46cb6-9b5a-4132-acd3-55a4fadcb248.jpg",
			md5sum: "d35b17ad7c9252f1a3daed176c211290",
		},
		{
			name:   "008-6eb93bee-5e94-4319-bb8d-2bbc1aabea80.png",
			path:   "artifacts/b/b/008-6eb93bee-5e94-4319-bb8d-2bbc1aabea80.png",
			md5sum: "4aad4f7d4cc3070d62180a259942697e",
		},
		{
			name:   "008-88d60783-5514-472f-9503-2904ab396e5f.jpg",
			path:   "artifacts/b/b/008-88d60783-5514-472f-9503-2904ab396e5f.jpg",
			md5sum: "f9446a65f8a96c32af58983be73d708d",
		},
		{
			name:   "008-8a1cf431-f1b3-4aa3-817c-5e55cc30b529.tif",
			path:   "artifacts/b/b/008-8a1cf431-f1b3-4aa3-817c-5e55cc30b529.tif",
			md5sum: "68b98eb3c13e9d5f4616810770826a01",
		},
		{
			name:   "008-b41fdceb-74c9-4a24-92e9-13f3ba086939.tif",
			path:   "artifacts/b/b/008-b41fdceb-74c9-4a24-92e9-13f3ba086939.tif",
			md5sum: "5aeae81585dea6451539a15ccefb1119",
		},
		{
			name:   "008-e819fae5-bd69-4e61-aace-b85f2ffd6d6e.png",
			path:   "artifacts/b/b/008-e819fae5-bd69-4e61-aace-b85f2ffd6d6e.png",
			md5sum: "e2a13e69ee33c76ae53caeecbe19db63",
		},
		{
			name:   "009-1ba73013-fecf-4a9c-bce9-d25a8d59f874.png",
			path:   "artifacts/b/b/009-1ba73013-fecf-4a9c-bce9-d25a8d59f874.png",
			md5sum: "86fb2342d7579f32748348e682188cac",
		},
		{
			name:   "009-372c2e12-85e1-4493-bb93-45baf2366faf.png",
			path:   "artifacts/b/b/009-372c2e12-85e1-4493-bb93-45baf2366faf.png",
			md5sum: "1c5e22ac6778dc4daeef4ba85fd7ac0a",
		},
		{
			name:   "009-8540ea54-68df-48e6-bb46-74fc8ade981d.tif",
			path:   "artifacts/b/b/009-8540ea54-68df-48e6-bb46-74fc8ade981d.tif",
			md5sum: "df62551b1c9731d0bb3ce40b1fa16974",
		},
		{
			name:   "009-aafd84b9-011e-4002-85be-703020f15ef8.jpg",
			path:   "artifacts/b/b/009-aafd84b9-011e-4002-85be-703020f15ef8.jpg",
			md5sum: "efd511565d133934a41ac158cd4b18ea",
		},
		{
			name:   "009-e0fd6627-17ce-4cc1-88e3-9eb8e6791cc0.jpg",
			path:   "artifacts/b/b/009-e0fd6627-17ce-4cc1-88e3-9eb8e6791cc0.jpg",
			md5sum: "26da9a3a581407ae189ff4f17310e4f4",
		},
		{
			name:   "009-fff05f6e-23a5-465d-bfbf-f6be9f906fe9.tif",
			path:   "artifacts/b/b/009-fff05f6e-23a5-465d-bfbf-f6be9f906fe9.tif",
			md5sum: "f2d7225652a46a49968daa33eb13774e",
		},
		{
			name:   "010-0360631f-6068-43c9-afa8-a2960e593309.tif",
			path:   "artifacts/b/b/010-0360631f-6068-43c9-afa8-a2960e593309.tif",
			md5sum: "86e107171a1db0b3bed62d7494bbcdb6",
		},
		{
			name:   "010-2b4fcf97-5985-480f-ad6b-24e724bada1d.png",
			path:   "artifacts/b/b/010-2b4fcf97-5985-480f-ad6b-24e724bada1d.png",
			md5sum: "216681600a43583a18a52e067cebf4fb",
		},
		{
			name:   "010-512abe3c-4b7a-4a48-829d-82dc795185ba.tif",
			path:   "artifacts/b/b/010-512abe3c-4b7a-4a48-829d-82dc795185ba.tif",
			md5sum: "4ea72ddfe90606bdfe2c9207ebe8ddc5",
		},
		{
			name:   "010-9b115de2-044e-4b0d-9bbf-d7d6bd8877c4.png",
			path:   "artifacts/b/b/010-9b115de2-044e-4b0d-9bbf-d7d6bd8877c4.png",
			md5sum: "dd7e4eed393b40d3453358d37bcdea37",
		},
		{
			name:   "010-adb8e778-5502-43e9-ace6-2ba2816f4ccf.jpg",
			path:   "artifacts/b/b/010-adb8e778-5502-43e9-ace6-2ba2816f4ccf.jpg",
			md5sum: "ba1bd87c8744faa647fbfccdf92e5957",
		},
		{
			name:   "010-bd68e013-baa5-48c2-864b-4ec951b3c235.jpg",
			path:   "artifacts/b/b/010-bd68e013-baa5-48c2-864b-4ec951b3c235.jpg",
			md5sum: "74120050621f6b77f3d70d4b0ec2439c",
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
			name:   "001-0d34c4ad-41ed-43fb-b8ae-1a660116a9a6.tif",
			path:   "artifacts/a/a/001-0d34c4ad-41ed-43fb-b8ae-1a660116a9a6.tif",
			md5sum: "94930845687e950296af035cb8637f9b",
		},
		{
			name:   "001-0d9e9496-cad5-4144-8285-593d1da5e396.png",
			path:   "artifacts/a/a/001-0d9e9496-cad5-4144-8285-593d1da5e396.png",
			md5sum: "7e0d87b875b8c574fff4cbc05fdc0f81",
		},
		{
			name:   "001-17ad62c2-08b9-4445-b75c-9764b530e2de.tif",
			path:   "artifacts/a/a/001-17ad62c2-08b9-4445-b75c-9764b530e2de.tif",
			md5sum: "e5bc221f7dafc7e4fbd0398f5d8aa93c",
		},
		{
			name:   "001-6f6faf97-8083-47b9-96a0-3f0406838abc.png",
			path:   "artifacts/a/a/001-6f6faf97-8083-47b9-96a0-3f0406838abc.png",
			md5sum: "8f2f76e86b31f39333edf6f40c8174fd",
		},
		{
			name:   "001-755695a2-d9e5-477f-bef2-5fab483d928b.jpg",
			path:   "artifacts/a/a/001-755695a2-d9e5-477f-bef2-5fab483d928b.jpg",
			md5sum: "9b899b6ac306d7f1e03ac62d3ecba207",
		},
		{
			name:   "001-9ccf7d26-2563-4a92-86ea-f4c7e79e4c69.jpg",
			path:   "artifacts/a/a/001-9ccf7d26-2563-4a92-86ea-f4c7e79e4c69.jpg",
			md5sum: "d9e212a57ea1189bf731556ea4c63cd8",
		},
		{
			name:   "002-05c5c532-e1d9-40e9-baa1-018dc26f7b87.png",
			path:   "artifacts/a/a/002-05c5c532-e1d9-40e9-baa1-018dc26f7b87.png",
			md5sum: "7030770d7e3bf1fa6c2b4945bcaec6ea",
		},
		{
			name:   "002-5fbc40ae-2006-478a-8047-ef0783e61967.jpg",
			path:   "artifacts/a/a/002-5fbc40ae-2006-478a-8047-ef0783e61967.jpg",
			md5sum: "b1681e0bf54ddf58ba5aa8ae5221fa15",
		},
		{
			name:   "002-76bcd5b7-33c5-40e9-8360-a9399a3f0309.tif",
			path:   "artifacts/a/a/002-76bcd5b7-33c5-40e9-8360-a9399a3f0309.tif",
			md5sum: "e33363f4c48ce09734b508c9349e3471",
		},
		{
			name:   "002-8e1a2964-1de9-4716-a5be-1d9e1979765e.png",
			path:   "artifacts/a/a/002-8e1a2964-1de9-4716-a5be-1d9e1979765e.png",
			md5sum: "49c2430894c84b6cda3feb2793c3764f",
		},
		{
			name:   "002-9ed8f982-db2b-444a-84f1-9a1cc06ea213.jpg",
			path:   "artifacts/a/a/002-9ed8f982-db2b-444a-84f1-9a1cc06ea213.jpg",
			md5sum: "e86d39b0746d6b3b6781cf8d5e12bc09",
		},
		{
			name:   "002-e61792ad-7b33-4030-8fe1-531a295a3827.tif",
			path:   "artifacts/a/a/002-e61792ad-7b33-4030-8fe1-531a295a3827.tif",
			md5sum: "4b01e1537895adc8fdef2fcd8183bd2f",
		},
		{
			name:   "003-31786aa2-e69a-43a5-93e8-6e8e58379850.jpg",
			path:   "artifacts/a/a/003-31786aa2-e69a-43a5-93e8-6e8e58379850.jpg",
			md5sum: "aa554f20575f620ef30f7b2c2e451a9e",
		},
		{
			name:   "003-4d7cd809-a46f-4de3-8718-26a9976835d7.tif",
			path:   "artifacts/a/a/003-4d7cd809-a46f-4de3-8718-26a9976835d7.tif",
			md5sum: "4039ad2bbbaf029df7ce88ba140d1fc7",
		},
		{
			name:   "003-773e329c-1cae-46c7-bf4d-eb7a87c6ef96.png",
			path:   "artifacts/a/a/003-773e329c-1cae-46c7-bf4d-eb7a87c6ef96.png",
			md5sum: "9bfe3f1482014223c44a0ca2d45c9851",
		},
		{
			name:   "003-c21b779c-596b-4746-b88b-117cdceac1e4.tif",
			path:   "artifacts/a/a/003-c21b779c-596b-4746-b88b-117cdceac1e4.tif",
			md5sum: "ae2dffe1946a048205623f58f8047699",
		},
		{
			name:   "003-ed5af8b4-5a81-409f-9b56-864f38f3f9e3.png",
			path:   "artifacts/a/a/003-ed5af8b4-5a81-409f-9b56-864f38f3f9e3.png",
			md5sum: "e2f257ccf736e4808a38e8a529bc20d3",
		},
		{
			name:   "003-fc596427-8316-4e58-b79f-a923076b283b.jpg",
			path:   "artifacts/a/a/003-fc596427-8316-4e58-b79f-a923076b283b.jpg",
			md5sum: "dd60b2631c82d46d69e43aa27683dd45",
		},
		{
			name:   "004-0a6cd19a-0b77-48e7-8062-03f3267a0238.jpg",
			path:   "artifacts/a/a/004-0a6cd19a-0b77-48e7-8062-03f3267a0238.jpg",
			md5sum: "ef949b5f4958d69578a4b000ea8bcab9",
		},
		{
			name:   "004-0aa80944-0b00-4ddb-9e16-4083abf7ff7a.jpg",
			path:   "artifacts/a/a/004-0aa80944-0b00-4ddb-9e16-4083abf7ff7a.jpg",
			md5sum: "d5f5b0d73da62445f74425684c973a26",
		},
		{
			name:   "004-9239d029-2c88-4b72-b9fc-9e33a91e2ff2.tif",
			path:   "artifacts/a/a/004-9239d029-2c88-4b72-b9fc-9e33a91e2ff2.tif",
			md5sum: "6b28dc32460743d30eb4b0c9dc32064b",
		},
		{
			name:   "004-9f677039-3691-46f7-bf1b-49ec65fce6b1.png",
			path:   "artifacts/a/a/004-9f677039-3691-46f7-bf1b-49ec65fce6b1.png",
			md5sum: "50b04bd870a05f8038ca259e3e0c42bc",
		},
		{
			name:   "004-a8816bc8-1fd2-4b60-a3e8-f52394e1a9fd.png",
			path:   "artifacts/a/a/004-a8816bc8-1fd2-4b60-a3e8-f52394e1a9fd.png",
			md5sum: "292e6f4cf96ca55e8d899d246cdabef7",
		},
		{
			name:   "004-dc867efd-af2a-4e23-9e50-923c1689714c.tif",
			path:   "artifacts/a/a/004-dc867efd-af2a-4e23-9e50-923c1689714c.tif",
			md5sum: "dba08c5a0248386d8ee01d135aa987a9",
		},
		{
			name:   "005-24040075-f733-4799-b353-56433a806d9f.tif",
			path:   "artifacts/a/a/005-24040075-f733-4799-b353-56433a806d9f.tif",
			md5sum: "f65a323c679bebe2fffb643354242eef",
		},
		{
			name:   "005-3dc2f3ec-7964-4070-8d54-230b75dff9b1.jpg",
			path:   "artifacts/a/a/005-3dc2f3ec-7964-4070-8d54-230b75dff9b1.jpg",
			md5sum: "31e3dec046759d710b230c716562caa1",
		},
		{
			name:   "005-8e897fcf-8376-4098-8789-8d2b7eae8084.png",
			path:   "artifacts/a/a/005-8e897fcf-8376-4098-8789-8d2b7eae8084.png",
			md5sum: "dc134f49e90b06de4bca5bd591b22580",
		},
		{
			name:   "005-b4480559-b551-4acb-abd0-5fe4622d8f8a.tif",
			path:   "artifacts/a/a/005-b4480559-b551-4acb-abd0-5fe4622d8f8a.tif",
			md5sum: "32ef05e82a5331fcb8c108a2311d7a93",
		},
		{
			name:   "005-f3aceaf7-c583-45e2-b376-000f2d1d8c89.jpg",
			path:   "artifacts/a/a/005-f3aceaf7-c583-45e2-b376-000f2d1d8c89.jpg",
			md5sum: "8ef1a0fb02c0358139ac8145be063c53",
		},
		{
			name:   "005-f69d7566-3f30-4998-9002-7a2ebccbb78e.png",
			path:   "artifacts/a/a/005-f69d7566-3f30-4998-9002-7a2ebccbb78e.png",
			md5sum: "187139fa2e2edc4478ec15b3020c6590",
		},
		{
			name:   "006-48e97bde-ba79-4540-9145-7783a86e3e02.jpg",
			path:   "artifacts/a/a/006-48e97bde-ba79-4540-9145-7783a86e3e02.jpg",
			md5sum: "db3e5866c52d2cbd4e6fbe14209d51a7",
		},
		{
			name:   "006-63106562-bd51-49dd-a9ba-f39757f5da8d.png",
			path:   "artifacts/a/a/006-63106562-bd51-49dd-a9ba-f39757f5da8d.png",
			md5sum: "b650de727860f075bb7db931f8150e0d",
		},
		{
			name:   "006-770fe9cc-12d8-4fd3-a41b-a32fa5bebf10.jpg",
			path:   "artifacts/a/a/006-770fe9cc-12d8-4fd3-a41b-a32fa5bebf10.jpg",
			md5sum: "36240d5b15454a46c40f70a2b6d46805",
		},
		{
			name:   "006-79f5183f-d8fa-44cb-87fd-1000349eda03.tif",
			path:   "artifacts/a/a/006-79f5183f-d8fa-44cb-87fd-1000349eda03.tif",
			md5sum: "a1b78d76ec28b50f604bab8a841a0bc3",
		},
		{
			name:   "006-910488af-fed2-4c2e-a615-8cbe67410e16.tif",
			path:   "artifacts/a/a/006-910488af-fed2-4c2e-a615-8cbe67410e16.tif",
			md5sum: "b9af47cc844822d6cc108ad6eecb6db5",
		},
		{
			name:   "006-b9e36e89-5431-422e-8326-e3857bc08b96.png",
			path:   "artifacts/a/a/006-b9e36e89-5431-422e-8326-e3857bc08b96.png",
			md5sum: "44ae4112a749d418077e3cac2c33e291",
		},
		{
			name:   "007-25d76443-0bfe-4971-93d4-1afeb72cba07.tif",
			path:   "artifacts/a/a/007-25d76443-0bfe-4971-93d4-1afeb72cba07.tif",
			md5sum: "57cbc7716917d6265376870750958a3d",
		},
		{
			name:   "007-5a21595c-b0b8-4c3a-94f3-b1fdab53bc57.tif",
			path:   "artifacts/a/a/007-5a21595c-b0b8-4c3a-94f3-b1fdab53bc57.tif",
			md5sum: "7ce154a96c041786946be65af335c71b",
		},
		{
			name:   "007-7730e806-5883-4311-b264-eac727d61504.png",
			path:   "artifacts/a/a/007-7730e806-5883-4311-b264-eac727d61504.png",
			md5sum: "d20d82035e7416d5acec23d4224271bc",
		},
		{
			name:   "007-b4731b77-55ef-4e7c-a844-c97ba873611b.png",
			path:   "artifacts/a/a/007-b4731b77-55ef-4e7c-a844-c97ba873611b.png",
			md5sum: "18105c380cffd8756fc535a080d5fb61",
		},
		{
			name:   "007-d54d8db6-8524-49b4-8481-23f27b4205fc.jpg",
			path:   "artifacts/a/a/007-d54d8db6-8524-49b4-8481-23f27b4205fc.jpg",
			md5sum: "2fa8d1930afbbf1518363f76eb4c9685",
		},
		{
			name:   "007-f57b5885-2ce4-4296-9448-763103a51513.jpg",
			path:   "artifacts/a/a/007-f57b5885-2ce4-4296-9448-763103a51513.jpg",
			md5sum: "7fa9b2e0efd498206a7582201c8f22b8",
		},
		{
			name:   "008-388cd1b6-180d-4ab8-ba94-84079435c8f8.png",
			path:   "artifacts/a/a/008-388cd1b6-180d-4ab8-ba94-84079435c8f8.png",
			md5sum: "dab046a085ac2fc8829009ccf8ceab34",
		},
		{
			name:   "008-8c1dc777-0c25-43ac-8b17-f1bf8ebc3d48.jpg",
			path:   "artifacts/a/a/008-8c1dc777-0c25-43ac-8b17-f1bf8ebc3d48.jpg",
			md5sum: "001d4d41ecabe1894f73c986067acf47",
		},
		{
			name:   "008-9551a46f-95d0-4125-b5e1-0d324187e48a.jpg",
			path:   "artifacts/a/a/008-9551a46f-95d0-4125-b5e1-0d324187e48a.jpg",
			md5sum: "fa3ae764cafd36adb4bcd9d9d2656620",
		},
		{
			name:   "008-ac6e1460-a079-4aa8-afc0-f9db0887e303.tif",
			path:   "artifacts/a/a/008-ac6e1460-a079-4aa8-afc0-f9db0887e303.tif",
			md5sum: "b60caca7355797f6b523aff0b387683f",
		},
		{
			name:   "008-af7403b5-04c9-4647-91a5-b1b1db90ac98.tif",
			path:   "artifacts/a/a/008-af7403b5-04c9-4647-91a5-b1b1db90ac98.tif",
			md5sum: "e6ed67a2a486c2b0159c85dc72146a68",
		},
		{
			name:   "008-e80e013f-5d6e-4152-89c0-ea321908fcc3.png",
			path:   "artifacts/a/a/008-e80e013f-5d6e-4152-89c0-ea321908fcc3.png",
			md5sum: "bc96a153e2198ef59516f791df1e1fa0",
		},
		{
			name:   "009-31aa9196-da7e-45c5-8ecf-b6872c9a6fbe.png",
			path:   "artifacts/a/a/009-31aa9196-da7e-45c5-8ecf-b6872c9a6fbe.png",
			md5sum: "fdb9a67c820ee17f86a410fd796ec5b9",
		},
		{
			name:   "009-5763e239-8c3b-435d-81ea-9cfe2db744d8.tif",
			path:   "artifacts/a/a/009-5763e239-8c3b-435d-81ea-9cfe2db744d8.tif",
			md5sum: "493ee38bce73bf60aa8f613350776a24",
		},
		{
			name:   "009-7d1810b9-ad56-4467-9ba0-80d36b9cad90.jpg",
			path:   "artifacts/a/a/009-7d1810b9-ad56-4467-9ba0-80d36b9cad90.jpg",
			md5sum: "18f86767325270f8a7fc26c5b8946774",
		},
		{
			name:   "009-7f6a3800-f123-418f-a4d5-7fbbbf61c7f2.jpg",
			path:   "artifacts/a/a/009-7f6a3800-f123-418f-a4d5-7fbbbf61c7f2.jpg",
			md5sum: "8db31fa5e190589ce88c4904f8d2d492",
		},
		{
			name:   "009-c54e4e9e-d848-479a-8a3e-60dcc5df85cc.png",
			path:   "artifacts/a/a/009-c54e4e9e-d848-479a-8a3e-60dcc5df85cc.png",
			md5sum: "a4732730d33b986eafd96e13ec5bd48c",
		},
		{
			name:   "009-ec10f514-f764-4dea-be42-73dcca158898.tif",
			path:   "artifacts/a/a/009-ec10f514-f764-4dea-be42-73dcca158898.tif",
			md5sum: "1165ab1741c59c831e1479daa3ef33e1",
		},
		{
			name:   "010-32e04124-2909-4fb4-802c-b0ce098c6d79.png",
			path:   "artifacts/a/a/010-32e04124-2909-4fb4-802c-b0ce098c6d79.png",
			md5sum: "68c706a8ddb8f43f71f3ab70ba67cdec",
		},
		{
			name:   "010-4089c12f-155a-4ca4-a568-89e647031e9a.jpg",
			path:   "artifacts/a/a/010-4089c12f-155a-4ca4-a568-89e647031e9a.jpg",
			md5sum: "f4ccc41e04cc16bfec6addc980b53970",
		},
		{
			name:   "010-4a9b928d-3ac5-4599-ae57-730c855a617f.jpg",
			path:   "artifacts/a/a/010-4a9b928d-3ac5-4599-ae57-730c855a617f.jpg",
			md5sum: "27501394b48a557aacad488cc6a0c7c6",
		},
		{
			name:   "010-5da64787-9780-4c4f-af9d-5a03ae9d4570.tif",
			path:   "artifacts/a/a/010-5da64787-9780-4c4f-af9d-5a03ae9d4570.tif",
			md5sum: "3ea8aff36c2fb464d4cffa7a7f3bcf24",
		},
		{
			name:   "010-7a7b93fe-2630-4706-8430-cb801032ef58.tif",
			path:   "artifacts/a/a/010-7a7b93fe-2630-4706-8430-cb801032ef58.tif",
			md5sum: "49ec45e103a8b84e86a323aba911855d",
		},
		{
			name:   "010-c432e094-3114-4aa0-b1c4-eacc25165e36.png",
			path:   "artifacts/a/a/010-c432e094-3114-4aa0-b1c4-eacc25165e36.png",
			md5sum: "5f1a2dfd4551f479c5097d510a89694b",
		},
	}
	got := readFiles([]string{"artifacts/a/a"})
	if !compareFileList(want, got) {
		t.Errorf("got: %s, want: %s\n", got, want)
	}
}

func TestRefreshAllFiles(t *testing.T) {
	oldTimestamp, _ := time.Parse(time.RFC822, "02 Jan 06 15:04 MST")
	var oldFiles Files = Files{
		{
			name:       "010-5da64787-9780-4c4f-af9d-5a03ae9d4570.tif",
			path:       "artifacts/a/a/010-5da64787-9780-4c4f-af9d-5a03ae9d4570.tif",
			md5sum:     "3ea8aff36c2fb464d4cffa7a7f3bcf24",
			lastPicked: oldTimestamp,
		},
		{
			name:   "010-7a7b93fe-2630-4706-8430-cb801032ef58.tif",
			path:   "artifacts/a/a/010-7a7b93fe-2630-4706-8430-cb801032ef58.tif",
			md5sum: "49ec45e103a8b84e86a323aba911855d",
		},
		{
			name:   "010-c432e094-3114-4aa0-b1c4-eacc25165e36.png",
			path:   "artifacts/a/a/010-c432e094-3114-4aa0-b1c4-eacc25165e36.png",
			md5sum: "5f1a2dfd4551f479c5097d510a89694b",
		},
	}
	var newFiles Files = Files{
		{
			name:   "010-4089c12f-155a-4ca4-a568-89e647031e9a.jpg",
			path:   "artifacts/a/a/010-4089c12f-155a-4ca4-a568-89e647031e9a.jpg",
			md5sum: "f4ccc41e04cc16bfec6addc980b53970",
		},
		{
			name:   "010-4a9b928d-3ac5-4599-ae57-730c855a617f.jpg",
			path:   "artifacts/a/a/010-4a9b928d-3ac5-4599-ae57-730c855a617f.jpg",
			md5sum: "27501394b48a557aacad488cc6a0c7c6",
		},
		{
			name:   "010-5da64787-9780-4c4f-af9d-5a03ae9d4570.tif",
			path:   "artifacts/a/a/010-5da64787-9780-4c4f-af9d-5a03ae9d4570.tif",
			md5sum: "3ea8aff36c2fb464d4cffa7a7f3bcf24",
		},
		{
			name:   "010-7a7b93fe-2630-4706-8430-cb801032ef58.tif",
			path:   "artifacts/a/a/010-7a7b93fe-2630-4706-8430-cb801032ef58.tif",
			md5sum: "49ec45e103a8b84e86a323aba911855d",
		},
	}
	var want Files = Files{
		{
			name:   "010-4089c12f-155a-4ca4-a568-89e647031e9a.jpg",
			path:   "artifacts/a/a/010-4089c12f-155a-4ca4-a568-89e647031e9a.jpg",
			md5sum: "f4ccc41e04cc16bfec6addc980b53970",
		},
		{
			name:   "010-4a9b928d-3ac5-4599-ae57-730c855a617f.jpg",
			path:   "artifacts/a/a/010-4a9b928d-3ac5-4599-ae57-730c855a617f.jpg",
			md5sum: "27501394b48a557aacad488cc6a0c7c6",
		},
		{
			name:       "010-5da64787-9780-4c4f-af9d-5a03ae9d4570.tif",
			path:       "artifacts/a/a/010-5da64787-9780-4c4f-af9d-5a03ae9d4570.tif",
			md5sum:     "3ea8aff36c2fb464d4cffa7a7f3bcf24",
			lastPicked: oldTimestamp,
		},
		{
			name:   "010-7a7b93fe-2630-4706-8430-cb801032ef58.tif",
			path:   "artifacts/a/a/010-7a7b93fe-2630-4706-8430-cb801032ef58.tif",
			md5sum: "49ec45e103a8b84e86a323aba911855d",
		},
	}
	var got Files = refreshAllFiles(oldFiles, newFiles)
	if !compareFileList(want, got) {
		t.Errorf("got: %s, want: %s\n", got, want)
	}
}
