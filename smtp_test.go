package hedwig

import (
	"context"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"testing"
)

func TestSMTPClient_SendMail(t *testing.T) {
	// Test run without error
	addr, ok := os.LookupEnv("TEST_SMTP_SERVER")
	if !ok {
		t.Skip()
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatal(err)
	}
	c := &SMTPClient{
		Address: addr,
		Auth:    smtp.PlainAuth("", os.Getenv("TEST_SMTP_USERNAME"), os.Getenv("TEST_SMTP_PASSWORD"), host),
	}
	ctx := context.TODO()
	m := &Mail{
		From: mail.Address{
			Name:    "送信者の名前",
			Address: os.Getenv("TEST_MAIL_FROM"),
		},
		To: []mail.Address{
			{Name: "受信者の名前", Address: os.Getenv("TEST_MAIL_TO")},
		},
		Subject:     "Generated from 銀河鉄道の夜",
		ContentType: ContentTypeHTML,
		Body: ` <h2>女の子は小さくほっと息をしながら言いました。</h2>
<p>またせっかくむいたそのきれいな砂を一つまみ、掌にひろげ、指でそっと、鷺のちぢめて降りて来る黒い脚をちぢめて、浮彫りのように幾本も幾本も、高く星ぞらに浮かんでいるのです。そのまっ黒な、松や楢の枝で、すっかりきれいに飾られた街を通って大通りへ出ていないかもしれないきっと出ている。双子のお星さまが野原へ遊びに来てくださいねそう言いながら、一枚の紙切れを渡しました。ところがその十字になったとこをはなして、三人の助手らしい人たちに囲まれた、小さな広場に出ました。ザネリもね、ずいぶん走ったけれども、いったいどんなことが、おっかさんは、なんにもひどいことないじゃないのジョバンニは靴をぬぎながら言いました。</p>
<ul>
<li>そしてのろしは高くそらにかかっているのでした。</li>
<li>すぐ前の席にいたわ女の子がこたえました。</li>
<li>川の遠くを飛んでいたのですか。</li>
</ul>
<h2>さがすと証拠もぞくぞく出ているのでした。</h2>
<p>こっち側の窓を見ましたら、向こうの席にすわったばかりの青い鋼の板のようなごうごうした声がきこえて来ましたので、ジョバンニも立って、ジョバンニは、ああ、そうだ、今夜ケンタウル祭だねえああ、ここは厚い立派な地層で、百二十万年ぐらい前のくるみだよ。わたしたちはもう、なんにもひどいことないじゃないのジョバンニは靴をぬぎながら言いました。なにがしあわせかわからない、そしてその一つの小さな星に見えるのでした。室中のひとたちは半分うしろの方からお持ちになったんですか。ああほんとうにどこまでも僕といっしょに歩いていたことは、二つにちぎってわたしました。</p>
<ol>
<li>おりるしたくをして言おうとしてしまいました。</li>
<li>ごとごと音をたてて流れているのでした。</li>
<li>ここでおりなけぁよかったなあ。</li>
</ol>`,
	}
	if err := c.SendMail(ctx, m); err != nil {
		t.Errorf("SMTPClient.SendMail() error = %v", err)
	}
}
