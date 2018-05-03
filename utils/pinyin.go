package utils

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"bytes"
	"io/ioutil"
	"os"
	"fmt"
	"bufio"
	"io"
	"strings"
	"strconv"
)
var (
	_GB2312_LETTER =[]rune("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbp" +
          "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbpbbbbbbbbbbbbbbbbbb" +
          "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" +
          "pbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" +
          "bbbbbbbbbbbbbbbbbbbbcccccccccccccccccccccccccccccc" +
          "ccccccccccccccccccccccccccccccccccczcccccccccccccc" +
          "ccccccccccccccccccccccccccccccccccccsccccccccccccc" +
          "cccccccccccccccccccccccccccccccccccccccccczccccccc" +
          "cccccccccccccccccccccccccccccccccccccccccccccccccc" +
          "cccddddddddddddddddddddddddddddddddddddddddddddddd" +
          "dddddddddddddddddddddzdddddddddddddddddddddddddddd" +
          "dddddddddddddddddddddddddddddddtdddddddddddddddddd" +
          "dddddddddddddddddddddddddddddddddddddeeeeeeeeeeeee" +
          "eeeeeeeeefffffffffffffffffffffffffffffffffffffffff" +
          "ffffffffffffffffffffffffffffffffffffffffffffffffff" +
          "fffffffffffffpffffffffffffffffffffgggggggggggggggg" +
          "ggggggggggggggggggghggggggggggggghgggggggggggggggg" +
          "gggggggggggggggggggggggggggggggggggggggggggggggggg" +
          "ggggggggggggggggggggggggggggggggggggggghhhhhhhhhhh" +
          "hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhmhhhhhhhhhhh" +
          "hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh" +
          "hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh" +
          "hhhhhhhhhhhhhhhhhhhhjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj" +
          "jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj" +
          "jjjjjjjjjjjjjjkjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj" +
          "jjjjjjjyjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj" +
          "jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj" +
          "jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj" +
          "jjjjjjjjjjjjjjjkkkgkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkh" +
          "kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk" +
          "kkkkkkkkkkkkkkklllllllllllllllllllllllllllllllllll" +
          "llllllllllllllllllllllllllllllllllllllllllllllllll" +
          "llllllllllllllllllllllllllllllllllllllllllllllllll" +
          "llllllllllllllllllllllllllllllllllllllllllllllllll" +
          "llllllllllllllllllllllllllllllllllllllllllllllllll" +
          "lllllllllllllmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm" +
          "mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm" +
          "mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm" +
          "mmmmmmmmmmmmmmnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn" +
          "nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnooooo" +
          "oooppppppppppppppppppppppppppppppppppppppppppppppp" +
          "pppppppppppppppppppppppppppppppppppppppppppppppppp" +
          "ppppppppppppppppppppppppbqqqqqqqqqqqqqqqqqqqqqqqqq" +
          "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq" +
          "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq" +
          "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrrrrrrrrrrrrrrrrrr" +
          "rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrsssssssss" +
          "ssssssssssssssssssssssssssssssssssssssssssssssssss" +
          "ssssssssssssssssssssssssssssssssssssssssssssssssss" +
          "ssssssssssssssssssssssssssssssssssssssssssssssssss" +
          "ssssssssssssssssssssssssssssssssssssssssssssssssss" +
          "sssssssssssssssssssssssssssssssssssssssssssssssssx" +
          "sssssssssssssssssssssssssssttttttttttttttttttttttt" +
          "tttttttttttttttttttttttttttttttttttttttttttttttttt" +
          "tttttttttttttttttttttttttttttttttttttttttttttttttt" +
          "tttttttttttttttttttttttttttttttttwwwwwwwwwwwwwwwww" +
          "wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww" +
          "wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww" +
          "wwwxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxsx" +
          "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
          "xxxxxxxxxxxxxxxxxxxxxjxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
          "xxxxxhxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxcxxxxxxxxx" +
          "xxxxxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyy" +
          "yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy" +
          "yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy" +
          "yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy" +
          "yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy" +
          "yyyyyyyyyyyyyyyyyyyyyyyyxyyyyyyyyyyyyyyyyyyyyyyyyy" +
          "yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzz" +
          "zzzzzzzzzzzzzzzzzzzzzczzzzzzzzzzzzzzzzzzzzzzzzzzzz" +
          "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz" +
          "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz" +
          "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz" +
          "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz" +
          "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz" +
          "zzzzz     cjwgnspgcgnesypbtyyzdxykygtdjnnjqmbsjzsc" +
          "yjsyyfpgkbzgylywjkgkljywkpjqhytwddzlsymrypywwcckzn" +
          "kyygttngjnykkzytcjnmcylqlypysfqrpzslwbtgkjfyxjwzlt" +
          "bncxjjjjtxdttsqzycdxxhgckbphffsswybgmxlpbylllhlxst" +
          "zmyjhsojnghdzqyklgjhsgqzhxqgkxzzwyscscjxyeyxadzpmd" +
          "ssmzjzqjyzcdjzwqjbyzbjgznzcpwhwxhqkmwfbpbydtjzzkxx" +
          "ylygxfptyjyyzpszlfchmqshgmxxsxjyqdcsbbqbefsjyhxwgz" +
          "kpylqbgldlcdtnmaeddkssngycsgxlyzaypnptsdkdylhgymyl" +
          "cxpycjndqjwxqxfyyfjlejpzrxccqwqqsbzkymgplbmjrqcfln" +
          "ymyqmtqyrbcjthztqfrxqhxmqjcjlyxgjmshzkbswyemyltxfs" +
          "ydsglycjqxsjnqbsctyhbftdcyjdjyyghqfsxwckqkxebptlpx" +
          "jzsrmebwhjlpjslyysmdxlclqkxlhxjrzjmfqhxhwywsbhtrxx" +
          "glhqhfnmgykldyxzpylggtmtcfpnjjzyljtyanjgbjplqgszyq" +
          "yaxbkysecjsznslyzhzxlzcghpxzhznytdsbcjkdlzayfmytle" +
          "bbgqyzkggldndnyskjshdlyxbcgyxypkdjmmzngmmclgezszxz" +
          "jfznmlzzthcsydbdllscddnlkjykjsycjlkwhqasdknhcsgaeh" +
          "daashtcplcpqybsdmpjlpzjoqlcdhjxysprchnwjnlhlyyqyhw" +
          "zptczgwwmzffjqqqqyxaclbhkdjxdgmmydqxzllsygxgkjrywz" +
          "wyclzmssjzldbydcpcxyhlxchyzjqsfqagmnyxpfrkssbjlyxy" +
          "syglnscmhcwwmnzjjlxxhchsyzsttxrycyxbyhcsmxjsznpwgp" +
          "xxtaybgajcxlysdccwzocwkccsbnhcpdyznfcyytyckxkybsqk" +
          "kytqqxfcwchcykelzqbsqyjqcclmthsywhmktlkjlycxwheqqh" +
          "tqkjpqsqscfymmdmgbwhwlgsllystlmlxpthmjhwljzyhzjxht" +
          "xjlhxrswlwzjcbxmhzqxsdzpsgfcsglsxymqshxpjxwmyqksmy" +
          "plrthbxftpmhyxlchlhlzylxgsssstclsldclrpbhzhxyyfhbb" +
          "gdmycnqqwlqhjjzywjzyejjdhpblqxtqkwhlchqxagtlxljxms" +
          "ljhtzkzjecxjcjnmfbycsfywybjzgnysdzsqyrsljpclpwxsdw" +
          "ejbjcbcnaytwgmpapclyqpclzxsbnmsggfnzjjbzsfzyndxhpl" +
          "qkzczwalsbccjxjyzgwkypsgxfzfcdkhjgxtlqfsgdslqwzkxt" +
          "mhsbgzmjzrglyjbpmlmsxlzjqzhzyjczydjwfmjklddpmjegxy" +
          "hylxhlqyqhkycwcjmyyxnatjhyccxzpcqlbzwwytwsqcmlpmyr" +
          "jcccxfpznzzljplxxyztzlgdltcklyrzzgqtkjhhgjljaxfgfj" +
          "zslcfdqzlclgjdjcsnzlljpjqdcclcjxmyzftsxgcgsbrzxjqq" +
          "ctzhgyqtjqqlzxjylylncyamcstylpdjbyregklzyzhlyszqlz" +
          "nwczcllwjqjjjkdgjzolbbzppglghtgzxyjhzmycnqsycyhbhg" +
          "xkamtxyxnbskyzzgjzlqjtfcjxdygjqjjpmgwgjjjpkqsbgbmm" +
          "cjssclpqpdxcdyykyfcjddyygywrhjrtgznyqldkljszzgzqzj" +
          "gdykshpzmtlcpwnjyfyzdjcnmwescyglbtzcgmssllyxqsxxbs" +
          "jsbbsgghfjlypmzjnlyywdqshzxtyywhmcyhywdbxbtlmsyyyf" +
          "sxjchtxxlhjhfssxzqhfzmzcztqcxzxrttdjhnnyzqqmtqdmmz" +
          " ytxmjgdxcdyzbffallztdltfxmxqzdngwqdbdczjdxbzgsqqd" +
          "djcmbkzffxmkdmdsyyszcmljdsynsprskmkmpcklgdbqtfzswt" +
          "fgglyplljzhgjjgypzltcsmcnbtjbqfkdhpyzgkpbbymtdssxt" +
          "bnpdkleycjnyddykzddhqhsdzsctarlltkzlgecllkjlqjaqnb" +
          "dkkghpjxzqksecshalqfmmgjnlyjbbtmlyzxdxjpldlpcqdhzy" +
          "cbzsczbzmsljflkrzjsnfrgjhxpdhyjybzgdlqcsezgxlblhyx" +
          "twmabchecmwyjyzlljjyhlgbdjlslygkdzpzxjyyzlwcxszfgw" +
          "yydlyhcljscmbjhblyzlycblydpdqysxqzbytdkyxlyycnrjmp" +
          "dqgklcljbcxbjddbblblczqrppxjcjlzcshltoljnmdddlngka" +
          "thqhjhykheznmshrphqqjchgmfprxhjgdychgklyrzqlcyqjnz" +
          "sqtkqjymszxwlcfqqqxyfggyptqwlmcrnfkkfsyylybmqammmy" +
          "xctpshcptxxzzsmphpshmclmldqfyqxszyjdjjzzhqpdszglst" +
          "jbckbxyqzjsgpsxqzqzrqtbdkwxzkhhgflbcsmdldgdzdblzyy" +
          "cxnncsybzbfglzzxswmsccmqnjqsbdqsjtxxmbltxcclzshzcx" +
          "rqjgjylxzfjphymzqqydfqjqlzznzjcdgzygztxmzysctlkpht" +
          "xhtlbjxjlxscdqxcbbtjfqzfsltjbtkqbxxjjljchczdbzjdcz" +
          "jdcprnpqcjpfczlclzxzdmxmphjsgzgszzqlylwtjpfsyaxmcj" +
          "btzyycwmytzsjjlqcqlwzmalbxyfbpnlsfhtgjwejjxxglljst" +
          "gshjqlzfkcgnndszfdeqfhbsaqtgylbxmmygszldydqmjjrgbj" +
          "tkgdhgkblqkbdmbylxwcxyttybkmrtjzxqjbhlmhmjjzmqasld" +
          "cyxyqdlqcafywyxqhz")

	_GBK_3 =[]rune("ksxsm sdqlybjjjgczbjfya jhphsyzgj   sn      xy  ng" +
"    lggllyjds yssgyqyd xjyydldwjjwbbftbxthhbczcrfm" +
"qwyfcwdzpyddwyxjajpsfnzyjxxxcxnnxxzzbpysyzhmzbqbzc" +
"ycbxqsbhhxgfmbhhgqcxsthlygymxalelccxzrcsd njjtzzcl" +
"jdtstbnxtyxsgkwyflhjqspxmxxdc lshxjbcfybyxhczbjyzl" +
"wlcz gtsmtzxpqglsjfzzlslhdzbwjncjysnycqrzcwybtyftw" +
"ecskdcbxhyzqyyxzcffzmjyxxsdcztbzjwszsxyrnygmdthjxs" +
"qqccsbxrytsyfbjzgclyzzbszyzqscjhzqydxlbpjllmqxtydz" +
"sqjtzplcgqtzwjbhcjdyfxjelbgxxmyjjqfzasyjnsydk jcjs" +
"zcbatdclnjqmwnqncllkbybzzsyhjqltwlccxthllzntylnzxd" +
"dtcenjyskkfksdkghwnlsjt jymrymzjgjmzgxykymsmjklfxm" +
"tghpfmqjsmtgjqdgyalcmzcsdjlxdffjc f  ffkgpkhrcjqcj" +
"dwjlfqdmlzbjjscgckdejcjdlzyckscclfcq czgpdqzjj hdd" +
"wgsjdkccctllpskghzzljlgjgjjtjjjzczmlzyjkxzyzmljkyw" +
"xmkjlkjgmclykjqlblkmdxwyxysllpsjqjqxyqfjtjdmxxllcr" +
"qyjb xgg pjygegdjgnjyjkhqfqzkhyghdgllsdjjxkyoxnzsx" +
"wwxdcskxxjyqscsqkjexsyzhydz ptqyzmtstzfsyldqagylcq" +
"lyyyhlrq ldhsssadsjbrszxsjyrcgqc hmmxzdyohycqgphhy" +
"nxrhgjlgwqwjhcstwasjpmmrdsztxyqpzxyhyqxtpbfyhhdwzb" +
"txhqeexzxxkstexgltxydn  hyktmzhxlplbmlsfhyyggbhyqt" +
"xwlqczydqdq gd lls zwjqwqajnytlxanzdecxzwwsgqqdyzt" +
"chyqzlxygzglydqtjtadyzzcwyzymhyhyjzwsxhzylyskqysbc" +
"yw  xjzgtyxqsyhxmchrwjpwxzlwjs sgnqbalzzmtjcjktsax" +
"ljhhgoxzcpdmhgtysjxhmrlxjkxhmqxctxwzbkhzccdytxqhlx" +
"hyx syydz znhxqyaygypdhdd pyzndltwxydpzjjcxmtlhbyn" +
"yymhzllhnmylllmdcppxmxdkycydltxchhznaclcclylzsxzjn" +
"zln lhyntkyjpychegttgqrgtgyhhlgcwyqkpyyyttttlhylly" +
"ttplkyzqqzdq  nmjzxyqmktfbjdjjdxbtqzgtsyflqgxblzfh" +
" zadpmjhlccyhdzfgydgcyxs hd d axxbpbyyaxcqffqyjxdl" +
"jjzl bjydyqszwjlzkcdtctbkdyzdqjnkknjgyeglfykasntch" +
"blwzbymjnygzyheyfjmctyfzjjhgck lxhdwxxjkyykssmwctq" +
"zlpbzdtwzxzag kwxl lspbclloqmmzslbczzkdcz xgqqdcyt" +
"zqwzqssfpktfqdcdshdtdwfhtdy jaqqkybdjyxtlj drqxxxa" +
"ydrjlklytwhllrllcxylbw z  zzhkhxksmdsyyjpzbsqlcxxn" +
"xwmdq gqmmczjgttybhyjbetpjxdqhkzbhfdxkawtwajldyjsf" +
"hblddqjncxfjhdfjjwzpkzypcyzynxff ydbzznytxzembsehx" +
"fzmbflzrsymzjrdjgxhjgjjnzzxhgxhymlpeyyxtgqshxssxmf" +
"mkcctxnypszhzptxwywxyysljsqxzdleelmcpjclxsqhfwwtff" +
"tnqjjjdxhwlyznflnkyyjldx hdynrjtywtrmdrqhwqcmfjdyz" +
"hmyyxjwzqtxtlmrspwwchjb xygcyyrrlmpymkszyjrmysntpl" +
"nbpyyxmykyngjzznlzhhanmpgwjdzmxxmllhgdzxyhxkrycjmf" +
"fxyhjfssqlxxndyca nmtcjcyprrnytyqym sxndlylyljnlxy" +
"shqmllyzljzxstyzsmcqynzlxbnnylrqtryyjzzhsytxcqgxzs" +
"shmkczyqhzjnbh qsnjnzybknlqhznswxkhjyybqlbfl p bkq" +
"zxsddjmessmlxxkwnmwwwydkzggtggxbjtdszxnxwmlptfxlcx" +
"jjljzxnwxlyhhlrwhsc ybyawjjcwqqjzzyjgxpltzftpakqpt" +
"lc  xtx hklefdleegqymsawhmljtwyqlyjeybqfnlyxrdsctg" +
"gxyyn kyqctlhjlmkkcgygllldzydhzwpjzkdyzzhyyfqytyzs" +
"ezzlymhjhtwyzlkyywzcskqqtdxwctyjklwqbdqyncs szjlkc" +
"dcdtlzzacqqzzddxyplxzbqjylzllqdzqgyjyjsyxnyyynyjxk" +
"xdazwrdljyyynjlxllhxjcykynqcclddnyyykyhhjcl pb qzz" +
"yjxj fzdnfpzhddwfmyypqjrssqzsqdgpzjwdsjdhzxwybp gp" +
"tmjthzsbgzmbjczwbbzmqcfmbdmcjxljbgjtz mqdyxjzyctyz" +
"tzxtgkmybbcljssqymscx jeglxszbqjjlyxlyctsxmcwfa kb" +
"qllljyxtyltxdphnhfqyzyes sdhwdjbsztfd czyqsyjdzjqp" +
"bs j fbkjbxtkqhmkwjjlhhyyyyywyycdypczyjzwdlfwxwzzj" +
"cxcdjzczlxjjtxbfwpxzptdzbccyhmlxbqlrtgrhqtlf mwwjx" +
"jwcysctzqhxwxkjybmpkbnzhqcdtyfxbyxcbhxpsxt m sxlhk" +
"mzxydhwxxshqhcyxglcsqypdh my ypyyykzljqtbqxmyhcwll" +
"cyl ewcdcmlggqktlxkgndgzyjjlyhqdtnchxwszjydnytcqcb" +
"hztbxwgwbxhmyqsycmqkaqyncs qhysqyshjgjcnxkzycxsbxx" +
"hyylstyxtymgcpmgcccccmztasgqzjlosqylstmqsqdzljqqyp" +
"lcycztcqqpbqjclpkhz yyxxdtddsjcxffllxmlwcjcxtspyxn" +
"dtjsjwxqqjskyylsjhaykxcyydmamdqmlmczncybzkkyflmcsc" +
"lhxrcjjgslnmtjzzygjddzjzk qgjyyxzxxqhheytmdsyyyqlf" +
" zzdywhscyqwdrxqjyazzzdywbjwhyqszywnp  azjbznbyzzy" +
"hnscpjmqcy zpnqtbzjkqqhngccxchbzkddnzhjdrlzlsjljyx" +
"ytbgtcsqmnjpjsrxcfjqhtpzsyjwbzzzlstbwwqsmmfdwjyzct" +
"bwzwqcslqgdhqsqlyzlgyxydcbtzkpj gm pnjkyjynhpwsnsz" +
"zxybyhyzjqjtllcjthgdxxqcbywbwzggqrqzssnpkydznxqxjm" +
"y dstzplthzwxwqtzenqzw ksscsjccgptcslccgllzxczqthn" +
"jgyqznmckcstjskbjygqjpldxrgzyxcxhgdnlzwjjctsbcjxbf" +
"zzpqdhjtywjynlzzpcjdsqjkdxyajyemmjtdljyryynhjbngzj" +
"kmjxltbsllrzylcscnxjllhyllqqqlxymswcxsljmc zlnsdwt" +
"jllggjxkyhbpdkmmscsgxjcsdybxdndqykjjtxdygmzzdzslo " +
"yjsjzdlbtxxxqqjzlbylwsjjyjtdzqqzzzzjlzcdzjhpl qplf" +
"fjzysj zfpfzksyjjhxttdxcysmmzcwbbjshfjxfqhyzfsjybx" +
"pzlhmbxhzxfywdab lktshxkxjjzthgxh jxkzxszzwhwtzzzs" +
"nxqzyawlcwxfxyyhxmyyswqmnlycyspjkhwcqhyljmzxhmcnzh" +
"hxcltjplxyjhdyylttxfszhyxxsjbjyayrmlckd yhlrlllsty" +
"zyyhscszqxkyqfpflk ntljmmtqyzwtlll s rbdmlqjbcc qy" +
"wxfzrzdmcyggzjm  mxyfdxc shxncsyjjmpafyfnhyzxyezy " +
"sdl zztxgfmyyysnbdnlhpfzdcyfssssn zzdgpafbdbzszbsg" +
"cyjlm  z yxqcyxzlckbrbrbzcycjzeeyfgzlyzsfrtkqsxdcm" +
"z  jl xscbykjbbrxllfqwjhyqylpzdxczybdhzrbjhwnjtjxl" +
"kcfssdqyjkzcwjl b  tzlltlqblcqqccdfpphczlyygjdgwcf" +
"czqyyyqyrqzslszfcqnwlhjcjjczkypzzbpdc   jgx gdz  f" +
"gpsysdfwwjzjyxyyjyhwpbygxrylybhkjksftzmmkhtyysyyzp" +
"yqydywmtjjrhl   tw  bjycfnmgjtysyzmsjyjhhqmyrszwtr" +
"tzsskx gqgsptgcznjjcxmxgzt ydjz lsdglhyqgggthszpyj" +
"hhgnygkggmdzylczlxqstgzslllmlcskbljzzsmmytpzsqjcj " +
" zxzzcpshkzsxcdfmwrllqxrfzlysdctmxjthjntnrtzfqyhqg" +
"llg   sjdjj tqjlnyhszxcgjzypfhdjspcczhjjjzjqdyb ss" +
"lyttmqtbhjqnnygjyrqyqmzgcjkpd gmyzhqllsllclmholzgd" +
"yyfzsljc zlylzqjeshnylljxgjxlyjyyyxnbzljsszcqqzjyl" +
"lzldj llzllbnyl hxxccqkyjxxxklkseccqkkkcgyyxywtqoh" +
"thxpyxx hcyeychbbjqcs szs lzylgezwmysx jqqsqyyycmd" +
"zywctjsycjkcddjlbdjjzqysqqxxhqjohdyxgmajpchcpljsmt" +
"xerxjqd pjdbsmsstktssmmtrzszmldj rn sqxqydyyzbdsln" +
"fgpzmdycwfdtmypqwytjzzqjjrjhqbhzpjhnxxyydyhhnmfcpb" +
"zpzzlzfmztzmyftskyjyjzhbzzygh pzcscsjssxfjgdyzyhzc" +
"whcsexfqzywklytmlymqpxxskqjpxzhmhqyjs cjlqwhmybdhy" +
"ylhlglcfytlxcjscpjskphjrtxteylssls yhxscznwtdwjslh" +
"tqdjhgydphcqfzljlzptynlmjllqyshhylqqzypbywrfy js y" +
"p yrhjnqtfwtwrchygmm yyhsmzhngcelqqmtcwcmpxjjfyysx" +
"ztybmstsyjdtjqtlhynpyqzlcxznzmylflwby jgsylymzctdw" +
"gszslmwzwwqzsayysssapxwcmgxhxdzyjgsjhygscyyxhbbzjk" +
"ssmalxycfygmqyjycxjlljgczgqjcczotyxmtthlwtgfzkpzcx" +
"kjycxctjcyh xsgckxzpsjpxhjwpjgsqxxsdmrszzyzwsykyzs" +
"hbcsplwsscjhjlchhylhfhhxjsx lnylsdhzxysxlwzyhcldyh" +
"zmdyspjtqznwqpsswctst zlmssmnyymjqjzwtyydchqlxkwbg" +
"qybkfc jdlzllyylszydwhxpsbcmljscgbhxlqrljxysdwxzsl" +
"df hlslymjljylyjcdrjlfsyjfnllcqyqfjy szlylmstdjcyh" +
"zllnwlxxygyygxxhhzzxczqzfnwpypkpypmlgxgg dxzzkzfbx" +
"xlzptytswhzyxhqhxxxywzyswdmzkxhzphgchj lfjxptzthly" +
"xcrhxshxkjxxzqdcqyl jlkhtxcwhjfwcfpqryqxyqy gpggsc" +
"sxngkchkzxhflxjbyzwtsxxncyjjmwzjqrhfqsyljzgynslgtc" +
"ybyxxwyhhxynsqymlywgyqbbzljlpsytjzhyzwlrorjkczjxxy" +
"xchdyxyxxjddsqfxyltsfxlmtyjmjjyyxltcxqzqhzlyyxzh n" +
"lrhxjcdyhlbrlmrllaxksllljlxxxlycry lccgjcmtlzllyzz" +
"pcw jyzeckzdqyqpcjcyzmbbcydcnltrmfgyqbsygmdqqzmkql" +
"pgtbqcjfkjcxbljmswmdt  ldlppbxcwkcbjczhkphyyhzkzmp" +
"jysylpnyyxdb")

_GBK_4 =[]rune(
"kxxmzjxsttdzxxbzyshjpfxpqbyljqkyzzzwl zgfwyctjxjpy" +
"yspmsmydyshqy zchmjmcagcfbbhplxtyqx djgxdhkxxnbhrm" +
"lnjsltsmrnlxqjyzlsqglbhdcgyqyyhwfjybbyjyjjdpqyapfx" +
"cgjscrssyz lbzjjjlgxzyxyxsqkxbxxgcxpld wetdwwcjmbt" +
"xchxyxxfxllj fwdpzsmylmwytcbcecblgdbqzqfjdjhymcxtx" +
"drmjwrh xcjzylqdyhlsrsywwzjymtllltqcjzbtckzcyqjzqa" +
"lmyhwwdxzxqdllqsgjfjljhjazdjgtkhsstcyjfpszlxzxrwgl" +
"dlzr lzqtgslllllyxxqgdzybphl x bpfd   hy jcc dmzpp" +
"z cyqxldozlwdwyythcqsccrsslfzfp qmbjxlmyfgjb m jwd" +
"n mmjtgbdzlp hsymjyl hdzjcctlcl ljcpddqdsznbgzxxcx" +
"qycbzxzfzfjsnttjyhtcmjxtmxspdsypzgmljtycbmdkycsz z" +
"yfyctgwhkyjxgyclndzscyzssdllqflqllxfdyhxggnywyllsd" +
"lbbjcyjzmlhl xyyytdlllb b bqjzmpclmjpgehbcqax hhhz" +
"chxyhjaxhlphjgpqqzgjjzzgzdqybzhhbwyffqdlzljxjpalxz" +
"daglgwqyxxxfmmsypfmxsyzyshdzkxsmmzzsdnzcfp ltzdnmx" +
"zymzmmxhhczjemxxksthwlsqlzllsjphlgzyhmxxhgzcjmhxtx" +
"fwkmwkdthmfzzydkmsclcmghsxpslcxyxmkxyah jzmcsnxyym" +
"mpmlgxmhlmlqmxtkzqyszjshyzjzybdqzwzqkdjlfmekzjpezs" +
"wjmzyltemznplplbpykkqzkeqlwayyplhhaq jkqclhyxxmlyc" +
"cyskg  lcnszkyzkcqzqljpmzhxlywqlnrydtykwszdxddntqd" +
"fqqmgseltthpwtxxlwydlzyzcqqpllkcc ylbqqczcljslzjxd" +
"dbzqdljxzqjyzqkzljcyqdypp pqykjyrpcbymxkllzllfqpyl" +
"llmsglcyrytmxyzfdzrysyztfmsmcl ywzgxzggsjsgkdtggzl" +
"ldzbzhyyzhzywxyzymsdbzyjgtsmtfxqyjssdgslnndlyzzlrx" +
"trznzxnqfmyzjzykbpnlypblnzz jhtzkgyzzrdznfgxskgjtt" +
"yllgzzbjzklplzylxyxbjfpnjzzxcdxzyxzggrs jksmzjlsjy" +
"wq yhqjxpjzt lsnshrnypzt wchklpszlcyysjylybbwzpdwg" +
"cyxckdzxsgzwwyqyytctdllxwkczkkcclgcqqdzlqcsfqchqhs" +
"fmqzlnbbshzdysjqplzcd cwjkjlpcmz jsqyzyhcpydsdzngq" +
"mbsflnffgfsm q lgqcyybkjsrjhzldcftlljgjhtxzcszztjg" +
"gkyoxblzppgtgyjdhz zzllqfzgqjzczbxbsxpxhyyclwdqjjx" +
"mfdfzhqqmqg yhtycrznqxgpdzcszcljbhbzcyzzppyzzsgyhc" +
"kpzjljnsc sllxb mstldfjmkdjslxlsz p pgjllydszgql l" +
"kyyhzttnt  tzzbsz ztlljtyyll llqyzqlbdzlslyyzyfszs" +
"nhnc   bbwsk rbc zm  gjmzlshtslzbl q xflyljqbzg st" +
"bmzjlxfnb xjztsfjmssnxlkbhsjxtnlzdntljjgzjyjczxygy" +
"hwrwqnztn fjszpzshzjfyrdjfcjzbfzqchzxfxsbzqlzsgyft" +
"zdcszxzjbqmszkjrhyjzckmjkhchgtxkjqalxbxfjtrtylxjhd" +
"tsjx j jjzmzlcqsbtxhqgxtxxhxftsdkfjhzxjfj  zcdlllt" +
"qsqzqwqxswtwgwbccgzllqzbclmqqtzhzxzxljfrmyzflxys x" +
"xjk xrmqdzdmmyxbsqbhgcmwfwtgmxlzpyytgzyccddyzxs g " +
"yjyznbgpzjcqswxcjrtfycgrhztxszzt cbfclsyxzlzqmzlmp" +
" lxzjxslbysmqhxxz rxsqzzzsslyflczjrcrxhhzxq dshjsj" +
"jhqcxjbcynsssrjbqlpxqpymlxzkyxlxcjlcycxxzzlxlll hr" +
"zzdxytyxcxff bpxdgygztcqwyltlswwsgzjmmgtjfsgzyafsm" +
"lpfcwbjcljmzlpjjlmdyyyfbygyzgyzyrqqhxy kxygy fsfsl" +
"nqhcfhccfxblplzyxxxkhhxshjzscxczwhhhplqalpqahxdlgg" +
"gdrndtpyqjjcljzljlhyhyqydhz zczywteyzxhsl jbdgwxpc" +
"  tjckllwkllcsstknzdnqnttlzsszyqkcgbhcrrychfpfyrwq" +
"pxxkdbbbqtzkznpcfxmqkcypzxehzkctcmxxmx nwwxjyhlstm" +
"csqdjcxctcnd p lccjlsblplqcdnndscjdpgwmrzclodansyz" +
"rdwjjdbcxwstszyljpxloclgpcjfzljyl c cnlckxtpzjwcyx" +
"wfzdknjcjlltqcbxnw xbxklylhzlqzllzxwjljjjgcmngjdzx" +
"txcxyxjjxsjtstp ghtxdfptffllxqpk fzflylybqjhzbmddb" +
"cycld tddqlyjjwqllcsjpyyclttjpycmgyxzhsztwqwrfzhjg" +
"azmrhcyy ptdlybyznbbxyxhzddnh msgbwfzzjcyxllrzcyxz" +
"lwjgcggnycpmzqzhfgtcjeaqcpjcs dczdwldfrypysccwbxgz" +
"mzztqscpxxjcjychcjwsnxxwjn mt mcdqdcllwnk zgglcczm" +
"lbqjqdsjzzghqywbzjlttdhhcchflsjyscgc zjbypbpdqkxwy" +
"yflxncwcxbmaykkjwzzzrxy yqjfljphhhytzqmhsgzqwbwjdy" +
"sqzxslzyymyszg x hysyscsyznlqyljxcxtlwdqzpcycyppnx" +
"fyrcmsmslxglgctlxzgz g tc dsllyxmtzalcpxjtjwtcyyjb" +
"lbzlqmylxpghdlssdhbdcsxhamlzpjmcnhjysygchskqmc lwj" +
"xsmocdrlyqzhjmyby lyetfjfrfksyxftwdsxxlysjslyxsnxy" +
"yxhahhjzxwmljcsqlkydztzsxfdxgzjksxybdpwnzwpczczeny" +
"cxqfjykbdmljqq lxslyxxylljdzbsmhpsttqqwlhogyblzzal" +
"xqlzerrqlstmypyxjjxqsjpbryxyjlxyqylthylymlkljt llh" +
"fzwkhljlhlj klj tlqxylmbtxchxcfxlhhhjbyzzkbxsdqc j" +
"zsyhzxfebcqwyyjqtzyqhqqzmwffhfrbntpcjlfzgppxdbbztg" +
" gchmfly xlxpqsywmngqlxjqjtcbhxspxlbyyjddhsjqyjxll" +
"dtkhhbfwdysqrnwldebzwcydljtmxmjsxyrwfymwrxxysztzzt" +
"ymldq xlyq jtscxwlprjwxhyphydnxhgmywytzcs tsdlwdcq" +
"pyclqyjwxwzzmylclmxcmzsqtzpjqblgxjzfljjytjnxmcxs c" +
"dl dyjdqcxsqyclzxzzxmxqrjhzjphfljlmlqnldxzlllfypny" +
"ysxcqqcmjzzhnpzmekmxkyqlxstxxhwdcwdzgyyfpjzdyzjzx " +
"rzjchrtlpyzbsjhxzypbdfgzzrytngxcqy b cckrjjbjerzgy" +
"  xknsjkljsjzljybzsqlbcktylccclpfyadzyqgk tsfc xdk" +
"dyxyfttyh  wtghrynjsbsnyjhkllslydxxwbcjsbbpjzjcjdz" +
"bfxxbrjlaygcsndcdszblpz dwsbxbcllxxlzdjzsjy lyxfff" +
"bhjjxgbygjpmmmpssdzjmtlyzjxswxtyledqpjmygqzjgdblqj" +
"wjqllsdgytqjczcjdzxqgsgjhqxnqlzbxsgzhcxy ljxyxydfq" +
"qjjfxdhctxjyrxysqtjxyebyyssyxjxncyzxfxmsyszxy schs" +
"hxzzzgzcgfjdltynpzgyjyztyqzpbxcbdztzc zyxxyhhsqxsh" +
"dhgqhjhgxwsztmmlhyxgcbtclzkkwjzrclekxtdbcykqqsayxc" +
"jxwwgsbhjyzs  csjkqcxswxfltynytpzc czjqtzwjqdzzzqz" +
"ljjxlsbhpyxxpsxshheztxfptjqyzzxhyaxncfzyyhxgnxmywx" +
"tcspdhhgymxmxqcxtsbcqsjyxxtyyly pclmmszmjzzllcogxz" +
"aajzyhjmzxhdxzsxzdzxleyjjzjbhzmzzzqtzpsxztdsxjjlny" +
"azhhyysrnqdthzhayjyjhdzjzlsw cltbzyecwcycrylcxnhzy" +
"dzydtrxxbzsxqhxjhhlxxlhdlqfdbsxfzzyychtyyjbhecjkgj" +
"fxhzjfxhwhdzfyapnpgnymshk mamnbyjtmxyjcthjbzyfcgty" +
"hwphftwzzezsbzegpbmtskftycmhbllhgpzjxzjgzjyxzsbbqs" +
"czzlzccstpgxmjsftcczjz djxcybzlfcjsyzfgszlybcwzzby" +
"zdzypswyjgxzbdsysxlgzybzfyxxxccxtzlsqyxzjqdcztdxzj" +
"jqcgxtdgscxzsyjjqcc ldqztqchqqjzyezwkjcfypqtynlmkc" +
"qzqzbqnyjddzqzxdpzjcdjstcjnxbcmsjqmjqwwjqnjnlllwqc" +
"qqdzpzydcydzcttf znztqzdtjlzbclltdsxkjzqdpzlzntjxz" +
"bcjltqjldgdbbjqdcjwynzyzcdwllxwlrxntqqczxkjld tdgl" +
" lajjkly kqll dz td ycggjyxdxfrskstqdenqmrkq  hgkd" +
"ldazfkypbggpzrebzzykyqspegjjglkqzzzslysywqzwfqzylz" +
"zlzhwcgkyp qgnpgblplrrjyxcccyyhsbzfybnyytgzxylxczw" +
"h zjzblfflgskhyjzeyjhlplllldzlyczblcybbxbcbpnnzc r" +
" sycgyy qzwtzdxtedcnzzzty hdynyjlxdjyqdjszwlsh lbc" +
"zpyzjyctdyntsyctszyyegdw ycxtscysmgzsccsdslccrqxyy" +
"elsm xztebblyylltqsyrxfkbxsychbjbwkgskhhjh xgnlycd" +
"lfyljgbxqxqqzzplnypxjyqymrbsyyhkxxstmxrczzywxyhymc" +
"l lzhqwqxdbxbzwzmldmyskfmklzcyqyczqxzlyyzmddz ftqp" +
"czcyypzhzllytztzxdtqcy ksccyyazjpcylzyjtfnyyynrs y" +
"lmmnxjsmyb sljqyldzdpqbzzblfndsqkczfywhgqmrdsxycyt" +
"xnq jpyjbfcjdyzfbrxejdgyqbsrmnfyyqpghyjdyzxgr htk " +
"leq zntsmpklbsgbpyszbydjzsstjzytxzphsszsbzczptqfzm" +
"yflypybbjgxzmxxdjmtsyskkbzxhjcelbsmjyjzcxt mljshrz" +
"zslxjqpyzxmkygxxjcljprmyygadyskqs dhrzkqxzyztcghyt" +
"lmljxybsyctbhjhjfcwzsxwwtkzlxqshlyjzjxe mplprcglt " +
"zztlnjcyjgdtclklpllqpjmzbapxyzlkktgdwczzbnzdtdyqzj" +
"yjgmctxltgcszlmlhbglk  njhdxphlfmkyd lgxdtwzfrjejz" +
"tzhydxykshwfzcqshknqqhtzhxmjdjskhxzjzbzzxympagjmst" +
"bxlskyynwrtsqlscbpspsgzwyhtlksssw hzzlyytnxjgmjszs" +
"xfwnlsoztxgxlsmmlbwldszylkqcqctmycfjbslxclzzclxxks" +
"bjqclhjpsqplsxxckslnhpsfqqytxy jzlqldtzqjzdyydjnzp" +
"d cdskjfsljhylzsqzlbtxxdgtqbdyazxdzhzjnhhqbyknxjjq" +
"czmlljzkspldsclbblzkleljlbq ycxjxgcnlcqplzlznjtzlx" +
"yxpxmyzxwyczyhzbtrblxlcczjadjlmmmsssmybhb kkbhrsxx" +
"jmxsdynzpelbbrhwghfchgm  klltsjyycqltskywyyhywxbxq" +
"ywbawykqldq tmtkhqcgdqktgpkxhcpthtwthkshthlxyzyyda" +
"spkyzpceqdltbdssegyjq xcwxssbz dfydlyjcls yzyexcyy" +
"sdwnzajgyhywtjdaxysrltdpsyxfnejdy lxllqzyqqhgjhzyc" +
"shwshczyjxllnxzjjn fxmfpycyawddhdmczlqzhzyztldywll" +
"hymmylmbwwkxydtyldjpyw xjwmllsafdllyflb   bqtzcqlj" +
"tfmbthydcqrddwr qnysnmzbyytbjhp ygtjahg tbstxkbtzb" +
"kldbeqqhqmjdyttxpgbktlgqxjjjcthxqdwjlwrfwqgwqhckry" +
"swgftgygbxsd wdfjxxxjzlpyyypayxhydqkxsaxyxgskqhykf" +
"dddpplcjlhqeewxksyykdbplfjtpkjltcyyhhjttpltzzcdlsh" +
"qkzjqyste eywyyzy xyysttjkllpwmcyhqgxyhcrmbxpllnqt" +
"jhyylfd fxzpsftljxxjbswyysksflxlpplbbblbsfxyzsylff" +
"fscjds tztryysyffsyzszbjtbctsbsdhrtjjbytcxyje xbne" +
"bjdsysykgsjzbxbytfzwgenhhhhzhhtfwgzstbgxklsty mtmb" +
"yxj skzscdyjrcwxzfhmymcxlzndtdh xdjggybfbnbpthfjaa" +
"xwfpxmyphdttcxzzpxrsywzdlybbjd qwqjpzypzjznjpzjlzt" +
" fysbttslmptzrtdxqsjehbzyj dhljsqmlhtxtjecxslzzspk" +
"tlzkqqyfs gywpcpqfhqhytqxzkrsg gsjczlptxcdyyzss qz" +
"slxlzmycbcqbzyxhbsxlzdltcdjtylzjyyzpzylltxjsjxhlbr" +
"ypxqzskswwwygyabbztqktgpyspxbjcmllxztbklgqkq lsktf" +
"xrdkbfpftbbrfeeqgypzsstlbtpszzsjdhlqlzpmsmmsxlqqnk" +
"nbrddnxxdhddjyyyfqgzlxsmjqgxytqlgpbqxcyzy drj gtdj" +
"yhqshtmjsbwplwhlzffny  gxqhpltbqpfbcwqdbygpnztbfzj" +
"gsdctjshxeawzzylltyybwjkxxghlfk djtmsz sqynzggswqs" +
"phtlsskmcl  yszqqxncjdqgzdlfnykljcjllzlmzjn   scht" +
"hxzlzjbbhqzwwycrdhlyqqjbeyfsjxwhsr  wjhwpslmssgztt" +
"ygyqqwr lalhmjtqjcmxqbjjzjxtyzkxbyqxbjxshzssfjlxmx" +
"  fghkzszggylcls rjyhslllmzxelgl xdjtbgyzbpktzhkzj" +
"yqsbctwwqjpqwxhgzgdyfljbyfdjf hsfmbyzhqgfwqsyfyjgp" +
"hzbyyzffwodjrlmftwlbzgycqxcdj ygzyyyyhy xdwegazyhx" +
"jlzythlrmgrxxzcl   ljjtjtbwjybjjbxjjtjteekhwslj lp" +
"sfyzpqqbdlqjjtyyqlyzkdksqj yyqzldqtgjj  js cmraqth" +
"tejmfctyhypkmhycwj cfhyyxwshctxrljhjshccyyyjltktty" +
"tmxgtcjtzaxyoczlylbszyw jytsjyhbyshfjlygjxxtmzyylt" +
"xxypzlxyjzyzyybnhmymdyylblhlsyygqllscxlxhdwkqgyshq" +
"ywljyyhzmsljljxcjjyy cbcpzjmylcqlnjqjlxyjmlzjqlycm" +
"hcfmmfpqqmfxlmcfqmm znfhjgtthkhchydxtmqzymyytyyyzz" +
"dcymzydlfmycqzwzz mabtbcmzzgdfycgcytt fwfdtzqssstx" +
"jhxytsxlywwkxexwznnqzjzjjccchyyxbzxzcyjtllcqxynjyc" +
"yycynzzqyyyewy czdcjyhympwpymlgkdldqqbchjxy       " +
"                                                  " +
"                 sypszsjczc     cqytsjljjt   ")



)

func UTF82GBK(s []byte)([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}


func ToFirstPyLetter(source string) string {
	if source!="" {
		var letters []rune
		runes:=[]rune(source)

		for _,r:=range runes {
			letters=append(letters,toLetter(r))
		}

      return string(letters)
	}

	return ""

}

func toLetter( r rune) rune {
	if r>=65&&r<=90 { //大写的字母
		return r-65+97
	}

	if r>=97&&r<=122 { //小写的字母
		return r
	}

	if r>1000 {
         return cnToLetter(r)
	}

	return r
}


func cnToLetter(r rune) rune{
    gbkdata,_:= UTF82GBK([]byte(string(r)))
	ch1:=int(gbkdata[0])
	ch2:=int(gbkdata[1])
	if ch1 <= 254 && ch1 >= 170 {
		//优先处理GB-2312汉字.
		if ch2 > 160 {
			//查找GB-2312
			no := (ch1-176)*94 + (ch2 - 160)
			return _GB2312_LETTER[no-1]
		}else {
			//查找GBK_4
			no:= (ch1 - 170) * 97 + (ch2 - 63);
			return _GBK_4[no-1]
		}

	}else if (ch1 <= 160 && ch1 >= 129) {
		//查找GBK_3
		no := (ch1 - 129) * 191 + (ch2 - 63);
		return _GBK_3[no-1]
	}

	return r
}
var _PINYIN_DICT =make(map[string]string)

func InitDict(dict string) {
	fi, err := os.Open(dict)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(a)
		arrays := strings.Split(line, ":")
		key := arrays[0]
		value := arrays[1]
		_PINYIN_DICT[strings.ToLower(key)] = value

	}
}

func ToPinyin(str string,sep string) string {
	textQuoted := strconv.QuoteToASCII(str)
	textUnquoted := textQuoted[1 : len(textQuoted)-1]
	sUnicodev := strings.Split(textUnquoted, "\\u")
	var content []string
	for _, v := range sUnicodev {
		py := _PINYIN_DICT[v]
		if py!="" {
			arrays := strings.Split(py, ",")
			//对于多音字，只取第一个音
			if len(arrays) > 1 {
				content=append(content,arrays[0])
			}else {
				content=append(content,py)
			}
		}
	}

	//返回拼音
	return strings.Join(content,sep)
}