package devatlas

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E", func() {

	It("should open DBs", func() {
		Expect(testDB.Meta).NotTo(BeNil())
		Expect(testDB.Meta.Ver).To(MatchRegexp(`\d+\.\d+`))

		Expect(testDB.Values).NotTo(BeEmpty())
		Expect(testDB.Properties).NotTo(BeEmpty())
		Expect(testDB.Regexp).NotTo(BeEmpty())
		Expect(testDB.Tree).NotTo(BeNil())
		Expect(testDB.Tree.Children).NotTo(BeEmpty())
		Expect(testDB.Tree.Data).To(BeEmpty())

		Expect(testDB.UAR).NotTo(BeNil())
	})

	It("should find attributes", func() {
		for _, tc := range testCases {
			attrs := testDB.Find(tc.ua)
			bin, err := json.Marshal(attrs)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(bin)).To(Equal(tc.json), "for '%s'", tc.ua)
		}
	})

	It("should normalize attributes", func() {
		attrs := testDB.Find(testCases[0].ua)
		Expect(attrs["vendor"]).To(Equal("Apple"))
		Expect(attrs["wmv"]).To(Equal(false))
		Expect(attrs["yearReleased"]).To(Equal(int(2007)))
	})

	It("should apply UAR transformations", func() {
		attrs := testDB.Find(`Mozilla/5.0 (iPod touch; CPU iPhone OS 7_0_4 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Mobile/11B554a`)
		Expect(attrs["vendor"]).To(Equal("Apple"))
		Expect(attrs["osVersion"]).To(Equal("7_0_4"))
	})

})

// TEST SUITE HOOK

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "github.com/bsm/devatlas")
}

var testDB *DB
var testCases = []struct{ ua, json string }{
	{`Mozilla/5.0 (iPhone; U; CPU like Mac OS X; en) AppleWebKit/420+ (KHTML, like Gecko) Version/3.0`,
		`{"3gp.aac.lc":true,"3gp.amr.nb":true,"3gp.amr.wb":false,"3gp.h263":true,"3gp.h264.level10":true,"3gp.h264.level10b":true,"3gp.h264.level11":true,"3gp.h264.level12":true,"3gp.h264.level13":true,"3gpp":false,"3gpp2":false,"aac":true,"aacInVideo":true,"aacLtpInVideo":false,"amr":false,"amrInVideo":false,"awbInVideo":false,"browserName":"Safari","browserRenderingEngine":"WebKit","browserVersion":"3.0","camera":"2.0","cldc":"-","cookieSupport":true,"csd":false,"css.animations":true,"css.columns":true,"css.transforms":true,"css.transitions":true,"developerPlatform":"1.1","developerPlatformVersion":"11.26","devicePixelRatio":1,"diagonalScreenSize":"3.5","displayColorDepth":24,"displayHeight":480,"displayPpi":165,"displayWidth":320,"drmOmaCombinedDelivery":false,"drmOmaForwardLock":false,"drmOmaSeparateDelivery":false,"edge":true,"flashCapable":false,"gprs":true,"h263Type0InVideo":false,"h263Type3InVideo":false,"hscsd":false,"hsdpa":false,"hspaEvolved":false,"html.audio":true,"html.canvas":true,"html.inlinesvg":true,"html.svg":true,"html.video":true,"https":true,"id":205202,"image.Gif87":true,"image.Gif89a":true,"image.Jpg":true,"image.Png":true,"isBrowser":false,"isChecker":false,"isDownloader":false,"isEReader":false,"isFeedReader":false,"isFilter":false,"isGamesConsole":false,"isMediaPlayer":false,"isMobilePhone":true,"isRobot":false,"isSetTopBox":false,"isSpam":false,"isTV":false,"isTablet":false,"jqm":true,"js.applicationCache":true,"js.deviceMotion":true,"js.deviceOrientation":true,"js.geoLocation":true,"js.indexedDB":false,"js.json":true,"js.localStorage":true,"js.modifyCss":true,"js.modifyDom":true,"js.querySelector":true,"js.sessionStorage":true,"js.supportBasicJavaScript":true,"js.supportConsoleLog":true,"js.supportEventListener":true,"js.supportEvents":true,"js.touchEvents":true,"js.webGl":false,"js.webSockets":true,"js.webSqlDatabase":true,"js.webWorkers":true,"js.xhr":true,"jsr118":false,"jsr139":false,"jsr30":false,"jsr37":false,"lte":false,"lteAdvanced":false,"manufacturer":"Apple","marketingName":"iPhone","markup.wml1":false,"markup.xhtmlBasic10":true,"markup.xhtmlMp10":true,"markup.xhtmlMp11":true,"markup.xhtmlMp12":true,"memoryLimitDownload":0,"memoryLimitEmbeddedMedia":0,"memoryLimitMarkup":0,"midiMonophonic":false,"midiPolyphonic":false,"midp":"-","mobileDevice":true,"model":"iPhone","mp3":true,"mp4.aac.lc":true,"mp4.h264.level11":true,"mp4.h264.level13":true,"mpeg4":true,"mpeg4InVideo":true,"nfc":false,"osAndroid":false,"osBada":false,"osLinux":false,"osName":"iOS","osOsx":true,"osProprietary":"True","osRim":false,"osSymbian":false,"osVersion":4,"osWebOs":false,"osWindows":false,"osWindowsMobile":false,"osWindowsPhone":false,"osWindowsRt":false,"osiOs":true,"primaryHardwareType":"Mobile Phone","qcelp":false,"qcelpInVideo":false,"stream.3gp.aac.lc":false,"stream.3gp.amr.nb":false,"stream.3gp.amr.wb":false,"stream.3gp.h263":false,"stream.3gp.h264.level10":false,"stream.3gp.h264.level10b":false,"stream.3gp.h264.level11":false,"stream.3gp.h264.level12":false,"stream.3gp.h264.level13":false,"stream.httpLiveStreaming":true,"stream.mp4.aac.lc":false,"stream.mp4.h264.level11":false,"stream.mp4.h264.level13":false,"supportsClientSide":true,"touchScreen":true,"umts":false,"uriSchemeSms":true,"uriSchemeSmsTo":false,"uriSchemeTel":true,"usableDisplayHeight":415,"usableDisplayWidth":320,"vCardDownload":false,"vendor":"Apple","wmv":false,"yearReleased":2007}`},
	{`BlackBerry9700/5.0.0.593 Profile/MIDP-2.1 Configuration/CLDC-1.1 VendorID/1`,
		`{"3gp.aac.lc":true,"3gp.amr.nb":true,"3gp.amr.wb":false,"3gp.h263":true,"3gp.h264.level10":true,"3gp.h264.level10b":true,"3gp.h264.level11":true,"3gp.h264.level12":true,"3gp.h264.level13":true,"aac":true,"amr":true,"browserName":"UP.Browser","browserRenderingEngine":"Mango","camera":"3.15","cldc":"1.1","cookieSupport":true,"diagonalScreenSize":"2.44","displayColorDepth":16,"displayHeight":360,"displayPpi":246,"displayWidth":480,"drmOmaCombinedDelivery":false,"drmOmaForwardLock":true,"drmOmaSeparateDelivery":false,"edge":true,"flashCapable":true,"gprs":true,"hsdpa":true,"https":true,"id":1669808,"image.Gif87":true,"image.Gif89a":true,"image.Jpg":true,"image.Png":true,"isBrowser":false,"isChecker":false,"isDownloader":false,"isEReader":false,"isFeedReader":false,"isFilter":false,"isGamesConsole":false,"isMediaPlayer":false,"isMobilePhone":true,"isRobot":false,"isSetTopBox":false,"isSpam":false,"isTV":false,"isTablet":false,"js.deviceMotion":false,"js.deviceOrientation":false,"js.indexedDB":false,"js.modifyCss":true,"js.modifyDom":true,"js.querySelector":true,"js.supportBasicJavaScript":true,"js.supportEventListener":true,"js.supportEvents":true,"js.webGl":false,"js.webSockets":false,"js.xhr":true,"jsr118":true,"jsr139":true,"jsr30":true,"jsr37":true,"lteAdvanced":false,"manufacturer":"RIM","marketingName":"Bold","markup.wml1":true,"markup.xhtmlBasic10":true,"markup.xhtmlMp10":true,"markup.xhtmlMp11":false,"markup.xhtmlMp12":false,"memoryLimitMarkup":32768,"midiMonophonic":true,"midiPolyphonic":true,"midp":"2.1","mobileDevice":true,"model":"BlackBerry 9700","mp3":true,"mp4.aac.lc":true,"mp4.h264.level11":true,"mp4.h264.level13":true,"nfc":false,"osAndroid":false,"osBada":false,"osLinux":false,"osName":"RIM","osOsx":false,"osRim":true,"osSymbian":false,"osVersion":"5.0.0.593","osWebOs":false,"osWindows":false,"osWindowsMobile":false,"osWindowsPhone":false,"osWindowsRt":false,"osiOs":false,"primaryHardwareType":"Mobile Phone","qcelp":false,"stream.3gp.aac.lc":true,"stream.3gp.amr.nb":true,"stream.3gp.amr.wb":false,"stream.3gp.h263":true,"stream.3gp.h264.level10":true,"stream.3gp.h264.level10b":true,"stream.3gp.h264.level11":true,"stream.3gp.h264.level12":true,"stream.3gp.h264.level13":true,"stream.mp4.aac.lc":true,"stream.mp4.h264.level11":true,"stream.mp4.h264.level13":true,"supportsClientSide":true,"touchScreen":false,"umts":true,"uriSchemeSms":false,"uriSchemeSmsTo":false,"uriSchemeTel":true,"usableDisplayHeight":348,"usableDisplayWidth":460,"vCardDownload":true,"vendor":"RIM","wmv":true,"yearReleased":2009}`},
	{`Opera/9.80 (Android 2.2; Opera Mobi/-2118645896; U; pl) Presto/2.7.60 Version/10.5`,
		`{"browserName":"Opera Mobile","browserRenderingEngine":"Presto","isBrowser":false,"isChecker":false,"isDownloader":false,"isFeedReader":false,"isFilter":false,"isRobot":false,"isSpam":false,"mobileDevice":true,"osName":"Android"}`},
	{`Mozilla/5.0 (iPad; CPU OS 7_0_4 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Mobile/11B554a [FBAN/FBIOS;FBAV/6.9.1;FBBV/1102303;FBDV/iPad3,4;FBMD/iPad;FBSN/iPhone OS;FBSV/7.0.4;FBSS/2; FBCR/;FBID/tablet;FBLC/en_US;FBOP/1],gzip(gfe)`,
		`{"3gp.aac.lc":true,"3gp.amr.nb":true,"3gp.amr.wb":false,"3gp.h263":false,"3gp.h264.level10":true,"3gp.h264.level10b":true,"3gp.h264.level11":true,"3gp.h264.level12":true,"3gp.h264.level13":true,"aac":true,"browserName":"Safari","browserRenderingEngine":"WebKit","cookieSupport":true,"css.animations":true,"css.columns":true,"css.transforms":true,"css.transitions":true,"devicePixelRatio":1,"diagonalScreenSize":"9.7","displayColorDepth":16,"displayHeight":1024,"displayPpi":132,"displayWidth":768,"flashCapable":false,"html.audio":true,"html.canvas":true,"html.svg":true,"html.video":true,"https":true,"id":1826129,"image.Gif87":true,"image.Gif89a":true,"image.Jpg":true,"image.Png":true,"isBrowser":false,"isChecker":false,"isDownloader":false,"isEReader":false,"isFeedReader":false,"isFilter":false,"isGamesConsole":false,"isMediaPlayer":false,"isMobilePhone":false,"isRobot":false,"isSetTopBox":false,"isSpam":false,"isTV":false,"isTablet":true,"jqm":true,"js.applicationCache":true,"js.deviceMotion":true,"js.deviceOrientation":false,"js.geoLocation":true,"js.indexedDB":false,"js.json":true,"js.localStorage":true,"js.modifyCss":true,"js.modifyDom":true,"js.querySelector":true,"js.sessionStorage":true,"js.supportBasicJavaScript":true,"js.supportConsoleLog":true,"js.supportEventListener":true,"js.supportEvents":true,"js.touchEvents":true,"js.webGl":false,"js.webSockets":true,"js.webSqlDatabase":true,"js.xhr":true,"lteAdvanced":false,"manufacturer":"Apple","marketingName":"iPad","markup.wml1":false,"markup.xhtmlBasic10":true,"markup.xhtmlMp10":true,"markup.xhtmlMp11":true,"markup.xhtmlMp12":true,"memoryLimitDownload":0,"memoryLimitEmbeddedMedia":0,"memoryLimitMarkup":0,"midiMonophonic":false,"midiPolyphonic":false,"mobileDevice":true,"model":"iPad","mp3":true,"mp4.aac.lc":true,"mp4.h264.level11":true,"mp4.h264.level13":true,"nfc":false,"osAndroid":false,"osBada":false,"osLinux":false,"osName":"iOS","osOsx":true,"osProprietary":"IPhone OS","osRim":false,"osSymbian":false,"osVersion":"7_0_4","osWebOs":false,"osWindows":false,"osWindowsMobile":false,"osWindowsPhone":false,"osWindowsRt":false,"osiOs":true,"primaryHardwareType":"Tablet","qcelpInVideo":false,"stream.3gp.aac.lc":false,"stream.3gp.amr.nb":false,"stream.3gp.amr.wb":false,"stream.3gp.h263":false,"stream.3gp.h264.level10":false,"stream.3gp.h264.level10b":false,"stream.3gp.h264.level11":false,"stream.3gp.h264.level12":false,"stream.3gp.h264.level13":false,"stream.httpLiveStreaming":true,"stream.mp4.aac.lc":false,"stream.mp4.h264.level11":false,"stream.mp4.h264.level13":false,"supportsClientSide":true,"touchScreen":true,"uriSchemeSms":false,"uriSchemeSmsTo":false,"uriSchemeTel":false,"usableDisplayHeight":1024,"usableDisplayWidth":768,"vendor":"Apple","wmv":false,"yearReleased":2010}`},
	{`Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_0 like Mac OS X; en-us) AppleWebKit/532.9 (KHTML, like Gecko) Version/4.0.5 Mobile/8A293 Safari/6531.22.7`,
		`{"3gp.aac.lc":true,"3gp.amr.nb":true,"3gp.amr.wb":false,"3gp.h263":true,"3gp.h264.level10":true,"3gp.h264.level10b":true,"3gp.h264.level11":true,"3gp.h264.level12":true,"3gp.h264.level13":true,"3gpp":false,"3gpp2":false,"aac":true,"aacInVideo":true,"aacLtpInVideo":false,"amr":false,"amrInVideo":false,"awbInVideo":false,"browserName":"Safari","browserRenderingEngine":"WebKit","browserVersion":"4.0.5","camera":"2.0","cldc":"-","cookieSupport":true,"csd":false,"css.animations":true,"css.columns":true,"css.transforms":true,"css.transitions":true,"developerPlatform":"1.1","developerPlatformVersion":"11.26","devicePixelRatio":1,"diagonalScreenSize":"3.5","displayColorDepth":24,"displayHeight":480,"displayPpi":165,"displayWidth":320,"drmOmaCombinedDelivery":false,"drmOmaForwardLock":false,"drmOmaSeparateDelivery":false,"edge":true,"flashCapable":false,"gprs":true,"h263Type0InVideo":false,"h263Type3InVideo":false,"hscsd":false,"hsdpa":false,"hspaEvolved":false,"html.audio":true,"html.canvas":true,"html.inlinesvg":true,"html.svg":true,"html.video":true,"https":true,"id":205202,"image.Gif87":true,"image.Gif89a":true,"image.Jpg":true,"image.Png":true,"isBrowser":false,"isChecker":false,"isDownloader":false,"isEReader":false,"isFeedReader":false,"isFilter":false,"isGamesConsole":false,"isMediaPlayer":false,"isMobilePhone":true,"isRobot":false,"isSetTopBox":false,"isSpam":false,"isTV":false,"isTablet":false,"jqm":true,"js.applicationCache":true,"js.deviceMotion":true,"js.deviceOrientation":true,"js.geoLocation":true,"js.indexedDB":false,"js.json":true,"js.localStorage":true,"js.modifyCss":true,"js.modifyDom":true,"js.querySelector":true,"js.sessionStorage":true,"js.supportBasicJavaScript":true,"js.supportConsoleLog":true,"js.supportEventListener":true,"js.supportEvents":true,"js.touchEvents":true,"js.webGl":false,"js.webSockets":true,"js.webSqlDatabase":true,"js.webWorkers":true,"js.xhr":true,"jsr118":false,"jsr139":false,"jsr30":false,"jsr37":false,"lte":false,"lteAdvanced":false,"manufacturer":"Apple","marketingName":"iPhone","markup.wml1":false,"markup.xhtmlBasic10":true,"markup.xhtmlMp10":true,"markup.xhtmlMp11":true,"markup.xhtmlMp12":true,"memoryLimitDownload":0,"memoryLimitEmbeddedMedia":0,"memoryLimitMarkup":0,"midiMonophonic":false,"midiPolyphonic":false,"midp":"-","mobileDevice":true,"model":"iPhone","mp3":true,"mp4.aac.lc":true,"mp4.h264.level11":true,"mp4.h264.level13":true,"mpeg4":true,"mpeg4InVideo":true,"nfc":false,"osAndroid":false,"osBada":false,"osLinux":false,"osName":"iOS","osOsx":true,"osProprietary":"True","osRim":false,"osSymbian":false,"osVersion":"4_0","osWebOs":false,"osWindows":false,"osWindowsMobile":false,"osWindowsPhone":false,"osWindowsRt":false,"osiOs":true,"primaryHardwareType":"Mobile Phone","qcelp":false,"qcelpInVideo":false,"stream.3gp.aac.lc":false,"stream.3gp.amr.nb":false,"stream.3gp.amr.wb":false,"stream.3gp.h263":false,"stream.3gp.h264.level10":false,"stream.3gp.h264.level10b":false,"stream.3gp.h264.level11":false,"stream.3gp.h264.level12":false,"stream.3gp.h264.level13":false,"stream.httpLiveStreaming":true,"stream.mp4.aac.lc":false,"stream.mp4.h264.level11":false,"stream.mp4.h264.level13":false,"supportsClientSide":true,"touchScreen":true,"umts":false,"uriSchemeSms":true,"uriSchemeSmsTo":false,"uriSchemeTel":true,"usableDisplayHeight":415,"usableDisplayWidth":320,"vCardDownload":false,"vendor":"Apple","wmv":false,"yearReleased":2007}`},
	{`Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:18.0) Gecko/20100101 Firefox/18.0`,
		`{"isBrowser":true,"isChecker":false,"isDownloader":false,"isFeedReader":false,"isFilter":false,"isRobot":false,"isSpam":false,"mobileDevice":false}`},
	{`Mozilla/5.0 (Linux; U; Android 2.2; de-de; Nexus One Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1`,
		`{"3gp.aac.lc":true,"3gp.amr.nb":true,"3gp.amr.wb":false,"3gp.h263":true,"3gp.h264.level10":true,"3gp.h264.level10b":true,"3gp.h264.level11":true,"3gp.h264.level12":true,"3gp.h264.level13":true,"aac":true,"amr":true,"browserName":"Android Browser","browserRenderingEngine":"WebKit","cldc":"1.1","cookieSupport":true,"css.animations":true,"css.columns":true,"css.transforms":true,"css.transitions":true,"developerPlatform":"Android","devicePixelRatio":"1.5","diagonalScreenSize":"3.82","displayColorDepth":24,"displayHeight":800,"displayPpi":244,"displayWidth":480,"drmOmaCombinedDelivery":false,"drmOmaForwardLock":true,"drmOmaSeparateDelivery":false,"edge":true,"flashCapable":true,"gprs":true,"hsdpa":true,"html.audio":true,"html.canvas":true,"html.inlinesvg":false,"html.svg":false,"html.video":true,"https":true,"id":1824544,"image.Gif87":true,"image.Gif89a":false,"image.Jpg":true,"image.Png":true,"isBrowser":false,"isChecker":false,"isDownloader":false,"isEReader":false,"isFeedReader":false,"isFilter":false,"isGamesConsole":false,"isMediaPlayer":false,"isMobilePhone":true,"isRobot":false,"isSetTopBox":false,"isSpam":false,"isTV":false,"isTablet":false,"jqm":true,"js.applicationCache":true,"js.deviceMotion":false,"js.deviceOrientation":false,"js.geoLocation":true,"js.indexedDB":false,"js.json":true,"js.localStorage":true,"js.modifyCss":true,"js.modifyDom":true,"js.querySelector":true,"js.sessionStorage":true,"js.supportBasicJavaScript":true,"js.supportConsoleLog":true,"js.supportEventListener":true,"js.supportEvents":true,"js.touchEvents":true,"js.webGl":false,"js.webSockets":false,"js.webSqlDatabase":true,"js.webWorkers":false,"js.xhr":true,"jsr118":false,"jsr139":false,"jsr30":false,"jsr37":false,"lteAdvanced":false,"manufacturer":"HTC","marketingName":"Passion","markup.wml1":false,"markup.xhtmlBasic10":true,"markup.xhtmlMp10":true,"markup.xhtmlMp11":true,"markup.xhtmlMp12":true,"memoryLimitDownload":0,"memoryLimitEmbeddedMedia":0,"memoryLimitMarkup":0,"midiMonophonic":true,"midiPolyphonic":true,"midp":"2.0","mobileDevice":true,"model":"Nexus One","mp3":true,"mp4.aac.lc":true,"mp4.h264.level11":true,"mp4.h264.level13":true,"nfc":false,"osAndroid":true,"osBada":false,"osLinux":true,"osName":"Android","osOsx":false,"osRim":false,"osSymbian":false,"osVersion":"2.2","osWebOs":false,"osWindows":false,"osWindowsMobile":false,"osWindowsPhone":false,"osWindowsRt":false,"osiOs":false,"primaryHardwareType":"Mobile Phone","qcelp":false,"stream.3gp.aac.lc":true,"stream.3gp.amr.nb":false,"stream.3gp.amr.wb":false,"stream.3gp.h263":true,"stream.3gp.h264.level10":true,"stream.3gp.h264.level10b":false,"stream.3gp.h264.level11":false,"stream.3gp.h264.level12":false,"stream.3gp.h264.level13":true,"stream.mp4.aac.lc":false,"stream.mp4.h264.level11":true,"stream.mp4.h264.level13":true,"supportsClientSide":true,"touchScreen":true,"umts":true,"uriSchemeSms":true,"uriSchemeSmsTo":true,"uriSchemeTel":true,"usableDisplayHeight":720,"usableDisplayWidth":380,"vCardDownload":false,"vendor":"HTC","wmv":false,"yearReleased":2010}`},
}

var _ = BeforeSuite(func() {
	var err error
	testDB, err = OpenFile("sample.json")
	Expect(err).NotTo(HaveOccurred())
})

// Benchmark

func BenchmarkFind_Simple(b *testing.B) {
	db, err := OpenFile("sample.json")
	if err != nil {
		b.Fatal("failed to open db:", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Find(testCases[5].ua)
	}
}

func BenchmarkFind_Hard(b *testing.B) {
	db, err := OpenFile("sample.json")
	if err != nil {
		b.Fatal("failed to open db:", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Find(testCases[3].ua)
	}
}
