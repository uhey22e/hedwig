package hedwig

import (
	"context"
	"net/mail"
	"os"
	"testing"
)

func TestSMTPClient_SendMail(t *testing.T) {
	addr, ok := os.LookupEnv("TEST_SMTP_SERVER")
	if !ok {
		t.Skip()
	}
	type args struct {
		ctx context.Context
		m   *Mail
	}
	tests := []struct {
		name    string
		c       *SMTPClient
		args    args
		wantErr bool
	}{
		{
			c: &SMTPClient{
				Address: addr,
			},
			args: args{
				ctx: context.TODO(),
				m: &Mail{
					From: mail.Address{
						Name:    "日本語の名前",
						Address: "japanese_name@example.com",
					},
					To: []mail.Address{
						{Address: "to@example.com"},
					},
					Subject:     "日本語のメール件名",
					ContentType: ContentTypeHTML,
					Body: `
<style>
html {
	font-family: sans-serif;
}
</style>
<h3>どこへ行ったのか、ぼんやりしてだまっていました。</h3>
<p>もうそこらが一ぺんに化石させて、こっちを見ていた席に、ぬれたように思いました。まあ、あの烏カムパネルラのとなりの席を指さしました。あんなにくるっとまわっていました。わたくしたちは神さまに召されているけやきの木のような音が聞こえてきました。けれどもだんだん気をつけていました。</p>
<ul>
<li>おや、あの河原は月夜だろうかそっちを見ていました。</li>
<li>僕こんな愉快な旅はしたことでもする。</li>
<li>ちゃんと小さな水晶のお宮だよ男の子が言いました。</li>
</ul>
<h3>それから元気よく口笛を吹きました。</h3>
<p>汽車が小さな小屋の前を通って行くのでした。ジョバンニはもう頭を引っ込めて地図を見てあわてたようについていて誰かの来るのを二人は見ました。わたしの大事なタダシはいまどんな歌をうたってやすむとき、いつも窓からぼんやり白く見えていたのです。ところが先生は早くもそれをもとめている。そしてそのこどもの肩のあたりが、どうも見たことあって僕あります。</p>
<ol>
<li>けれども見つからないんだからお母さん。</li>
<li>こいつはすこしもいたようではありませんか。</li>
<li>右手の低い丘の上にかかえていました。</li>
</ol>`,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SendMail(tt.args.ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("SMTPClient.SendMail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
