package locale

import (
	"regexp"
)

// https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes

var alphabets = map[string]string{
	// aa,aar,Afar
	"aa": `ABCDEFGHIJKLMNOPQRSTUVWXYZ`,
	// ab,abk,Abkhazian
	"ab": `АБВГӶДЕЖЗӠИКҚҞЛМНОПԤРСТҬУФХҲЦҴЧҶҼҾШЫҨЏЬӘ`,
	// ae,ave,Avestan
	"ae": "𐬀𐬁𐬂𐬃𐬄𐬅𐬆𐬇𐬈𐬉𐬊𐬋𐬌𐬍𐬎𐬏𐬐𐬑𐬒𐬓𐬔𐬕𐬖𐬗𐬘𐬙𐬚𐬛𐬜𐬝𐬞𐬟𐬠𐬡𐬢𐬣𐬤𐬥𐬦𐬧𐬨𐬩𐬪𐬫𐬬𐬭𐬯𐬰𐬱𐬲𐬳𐬴𐬵",
	// af,afr,Afrikaans
	"af": `AÁBCDEÉÈÊËFGHIÍÎÏJKLMNOÓÔÖPQRSTUÚÛÜVWXYÝZ`,
	// ak,aka,Akan
	"ak": `ABDEƐFGHIKLMNOƆPRSTUWY`,
	// am,amh,Amharic (https://www.amharicalphabet.com/)
	"am": `ሀ ሁ ሂ ሃ ሄ ህ ሆ ለ ሉ ሊ ላ ሌ ል ሎ ሐ ሑ ሒ ሓ ሔ ሕ ሖ መ ሙ ሚ ማ ሜ ም ሞ
	ሠ ሡ ሢ ሣ ሤ ሥ ሦ ረ ሩ ሪ ራ ሬ ር ሮ ሰ ሱ ሲ ሳ ሴ ስ ሶ ሸ ሹ ሺ ሻ ሼ ሽ ሾ ቀ ቁ ቂ ቃ ቄ ቅ ቆ
	ቐ ቑ ቒ ቓ ቔ ቕ ቖ በ ቡ ቢ ባ ቤ ብ ቦ ቨ ቩ ቪ ቫ ቬ ቭ ቮ ተ ቱ ቲ ታ ቴ ት ቶ ቸ ቹ ቺ ቻ ቼ ች ቾ
	ኀ ኁ ኂ ኃ ኄ ኅ ኆ ነ ኑ ኒ ና ኔ ን ኖ ኘ ኙ ኚ ኛ ኜ ኝ ኞ አ ኡ ኢ ኣ ኤ እ ኦ ከ ኩ ኪ ካ ኬ ክ ኮ
	ኸ ኹ ኺ ኻ ኼ ኽ ኾ ወ ዉ ዊ ዋ ዌ ው ዎ ዐ ዑ ዒ ዓ ዔ ዕ ዖ ዘ ዙ ዚ ዛ ዜ ዝ ዞ ዠ ዡ ዢ ዣ ዤ ዥ ዦ
	የ ዩ ዪ ያ ዬ ይ ዮ ደ ዱ ዲ ዳ ዴ ድ ዶ ዸ ዹ ዺ ዻ ዼ ዽ ዾ ጀ ጁ ጂ ጃ ጄ ጅ ጆ ገ ጉ ጊ ጋ ጌ ግ ጎ
	ጘ ጙ ጚ ጛ ጜ ጝ ጞ ጠ ጡ ጢ ጣ ጤ ጥ ጦ ጨ ጩ ጪ ጫ ጬ ጭ ጮ ጰ ጱ ጲ ጳ ጴ ጵ ጶ ጸ ጹ ጺ ጻ ጼ ጽ ጾ
	ፀ ፁ ፂ ፃ ፄ ፅ ፆ ፈ ፉ ፊ ፋ ፌ ፍ ፎ ፐ ፑ ፒ ፓ ፔ ፕ ፖ`,
	// an,arg,Aragonese
	"an": `ABCDEFGHIJKLMNOPQRSTUVWXYZ`,
	// ar,ara,Arabic
	"ar": `ؠ ء آ أ ؤ إ ئ ا ب ة ت ث ج ح خ د ذ ر ز س ش ص ض ط ظ ع غ
	ػ ؼ ؽ ؾ ؿ ف ق ك ل م ه و ى ي ٮ ٯ ٱ ٲ ٳ ٴ ٵ ٶ ٷ ٸ ٹ ٺ ٻ ټ ٽ پ ٿ ڀ ځ ڂ
	ڃ ڄ څ چ ڇ ڈ ډ ڊ ڋ ڌ ڍ ڎ ڏ ڐ ڑ ڒ ړ ڔ ڕ ږ ڗ ژ ڙ ښ ڛ ڜ ڝ ڞ ڟ ڠ ڡ ڢ ڣ ڤ
	ڥ ڦ ڧ ڨ ک ڪ ګ ڬ ڭ ڮ گ ڰ ڱ ڲ ڳ ڴ ڵ ڶ ڷ ڸ ڹ ں ڻ ڼ ڽ ھ ڿ ۀ ہ ۂ ۃ ۄ ۅ ۆ
	ۇ ۈ ۉ ۊ ۋ ی ۍ ێ ۏ ې ۑ ے ۓ ە ۥ ۦ ۮ ۯ ۺ ۻ ۼ ۿ`,
	// as,asm,Assamese
	// https://en.wikipedia.org/wiki/Bengali%E2%80%93Assamese_script
	// https://www.unicode.org/charts/PDF/U0980.pdf
	// https://en.wikipedia.org/wiki/Assamese_alphabet
	"as": ``,
	// av,ava,Avaric
	"av": `АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯӀ`,
	// ay,aym,Aymara (https://aymara.org/webarchives/www2000/english/qillqa.html)
	"ay": `AÄCHIÏJKLMNÑPQRSTUÜWXY'`,
	// az,aze,Azerbaijani
	"az": `ABCÇDEƏFGĞHXIİJKQLMNOÖPRSŞTUÜVYZ`,
	////////////////////////////////////////////////////////////////////////
	// Bashkir,ba,bak
	// be,bel,Belarusian
	"be": `АБВГДЕЁЖЗІЙКЛМНОПРСТУЎФХЦЧШЫЬЭЮЯ`,
	// bg,bul,Bulgarian
	"bg": `АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЬЮЯ`,
	// Bihari languages,bh,bih
	// Bislama,bi,bis
	// Bambara,bm,bam
	// Bengali,bn,ben
	// Tibetan,bo,bod
	// Breton,br,bre
	// Bosnian,bs,bos
	////////////////////////////////////////////////////////////////////////
	// ,cnr,Montenegrin
	"cnr": `ABCČĆDĐEFGHIJKLMNOPRSŠŚTUVZŽŹАБВГДЂЕЖЗЗ́ИЈКЛЉМНЊОПРСС́ТЋУФХЦЧЏШ`,
	// cs,ces,Czech
	"cs": `AÁBCČDĎEÉĚFGHIÍJKLMNŇOÓPQRŘSŠTŤUÚŮVWXYÝZŽ`,
	// da,dan,Danish
	"da": `AÁBCDEÉFGHIÍJKLMNOÓPQRSTUÚVWXYÝZÆǼØǾÅ`,
	// de,deu,German
	"de": `ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜẞ`,
	// en,eng,English
	"en": `ABCDEFGHIJKLMNOPQRSTUVWXYZ`,
	// fo,fao,Faroese
	"fo": `AÁBDÐEFGHIÍJKLMNOÓPRSTUÚVYÝÆØ`,
	// is,isl,Icelandic
	"is": `AÁBDÐEÉFGHIÍJKLMNOÓPRSTUÚVXYÝÞÆÖ`,
	// mk,mkd,Macedonian
	"mk": `АБВГДЃЕЖЗЅИЈКЛЉМНЊОПРСТЌУФХЦЧЏШ`,
	// nl,nld,Dutch, Flemish
	"nl": `AÁÄBCDEÉËFGHIÍÏJJ́KLMNOÓÖPQRSTUÚÜVWXYÝZ`, // what about "Ÿ" ?
	// no,nor,Norwegian
	"no": `AÁÂĀBCDEÉÈÊĒFGHIÍĪJKLMNOÓÒÔŌPQRSTUVWXYÝZÆØÅ`, // is "Ç" used in loanwords? should it be included?
	// pl,pol,Polish
	"pl": `AĄBCĆDEĘFGHIJKLŁMNŃOÓPQRSŚTUVWXYZŹŻ`,
	// ru,rus,Russian
	"ru": `АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ`,
	// sh, ... ("sh" is deprecated, use "hbs")
	"sh": `ABCČĆDDžĐEFGHIJKLLMNNOPRSŠTUVZŽ`,
	// sk,slk,Slovak
	"sk": `AÁÄBCČDĎEÉFGHIÍJKLĹĽMNŇOÓÔPQRŔSŠTŤUÚVWXYÝZŽ`,
	// sl,slv,Slovenian
	"sl": `ABCČDEFGHIJKLMNOPRSŠTUVZŽ`,
	// sr,srp,Serbian
	"sr": `АБВГДЂЕЖЗИЈКЛЉМНЊОПРСТЋУФХЦЧЏШ`,
	// sv,swe,Swedish
	"sv": `ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ`,
	// uk,ukr,Ukrainian
	"uk": `АБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЬЮЯ`,
	// ,wen,sorbian languages
	"wen": `ABCČĆDEĚFGHIJKŁLMNŃOÓPRŘŔSŠŚTUWYZŽŹ`,
}

// Albanian,sq,sqi
// Armenian,hy,hye

// Basque,eu,eus
// Burmese,my,mya
// Catalan, Valencian,ca,cat
// Chamorro,ch,cha
// Chechen,ce,che
// Chichewa, Chewa, Nyanja,ny,nya
// Chinese,zh,zho
// Chuvash,cv,chv
// Cornish,kw,cor
// Corsican,co,cos
// Cree,cr,cre
// Croatian,hr,hrv
// Divehi, Dhivehi, Maldivian,dv,div
// Dzongkha,dz,dzo
// Esperanto,eo,epo
// Estonian,et,est
// Ewe,ee,ewe
// Fijian,fj,fij
// Finnish,fi,fin
// French,fr,fra
// Fulah,ff,ful
// Galician,gl,glg
// Georgian,ka,kat
// Greek, Modern (1453–),el,ell
// Guarani,gn,grn
// Gujarati,gu,guj
// Haitian, Haitian Creole,ht,hat
// Hausa (Hausa),ha,hau
// Hebrew,he,heb
// Herero,hz,her
// Hindi,hi,hin
// Hiri Motu,ho,hmo
// Hungarian,hu,hun
// Interlingua (International Auxiliary Language Association),ia,ina
// Indonesian,id,ind
// Interlingue, Occidental,ie,ile
// Irish,ga,gle
// Igbo,ig,ibo
// Inupiaq,ik,ipk
// Ido,io,ido
// Italian,it,ita
// Inuktitut,iu,iku
// Japanese,ja,jpn
// Javanese,jv,jav
// Kalaallisut, Greenlandic,kl,kal
// Kannada,kn,kan
// Kanuri,kr,kau
// Kashmiri,ks,kas
// Kazakh,kk,kaz
// Central Khmer,km,khm
// Kikuyu, Gikuyu,ki,kik
// Kinyarwanda,rw,kin
// Kirghiz, Kyrgyz,ky,kir
// Komi,kv,kom
// Kongo,kg,kon
// Korean,ko,kor
// Kurdish,ku,kur
// Kuanyama, Kwanyama,kj,kua
// Latin,la,lat
// Luxembourgish, Letzeburgesch,lb,ltz
// Ganda,lg,lug
// Limburgan, Limburger, Limburgish,li,lim
// Lingala,ln,lin
// Lao,lo,lao
// Lithuanian,lt,lit
// Luba-Katanga,lu,lub
// Latvian,lv,lav
// Manx,gv,glv
// Malagasy,mg,mlg
// Malay,ms,msa
// Malayalam,ml,mal
// Maltese,mt,mlt
// Maori,mi,mri
// Marathi,mr,mar
// Marshallese,mh,mah
// Mongolian,mn,mon
// Nauru,na,nau
// Navajo, Navaho,nv,nav
// North Ndebele,nd,nde
// Nepali,ne,nep
// Ndonga,ng,ndo
// Norwegian Bokmål,nb,nob
// Norwegian Nynorsk,nn,nno
// Sichuan Yi, Nuosu,ii,iii
// South Ndebele,nr,nbl
// Occitan,oc,oci
// Ojibwa,oj,oji
// Church Slavic, Old Slavonic, Church Slavonic, Old Bulgarian, Old Church Slavonic,cu,chu
// Oromo,om,orm
// Oriya,or,ori
// Ossetian, Ossetic,os,oss
// Punjabi, Panjabi,pa,pan
// Pali,pi,pli
// Persian,fa,fas
// Pashto, Pushto,ps,pus
// Portuguese,pt,por
// Quechua,qu,que
// Romansh,rm,roh
// Rundi,rn,run
// Romanian, Moldavian, Moldovan,ro,ron
// Sanskrit,sa,san
// Sardinian,sc,srd
// Sindhi,sd,snd
// Northern Sami,se,sme
// Samoan,sm,smo
// Sango,sg,sag
// Gaelic, Scottish Gaelic,gd,gla
// Shona,sn,sna
// Sinhala, Sinhalese,si,sin
// Somali,so,som
// Southern Sotho,st,sot
// Spanish, Castilian,es,spa
// Sundanese,su,sun
// Swahili,sw,swa
// Swati,ss,ssw
// Tamil,ta,tam
// Telugu,te,tel
// Tajik,tg,tgk
// Thai,th,tha
// Tigrinya,ti,tir
// Turkmen,tk,tuk
// Tagalog,tl,tgl
// Tswana,tn,tsn
// Tonga (Tonga Islands),to,ton
// Turkish,tr,tur
// Tsonga,ts,tso
// Tatar,tt,tat
// Twi,tw,twi
// Tahitian,ty,tah
// Uighur, Uyghur,ug,uig
// Urdu,ur,urd
// Uzbek,uz,uzb
// Venda,ve,ven
// Vietnamese,vi,vie
// Volapük,vo,vol
// Walloon,wa,wln
// Welsh,cy,cym
// Wolof,wo,wol
// Western Frisian,fy,fry
// Xhosa,xh,xho
// Yiddish,yi,yid
// Yoruba,yo,yor
// Zhuang, Chuang,za,zha
// Zulu,zu,zul
