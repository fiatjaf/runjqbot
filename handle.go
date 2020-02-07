package main

import (
	"encoding/json"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func handle(upd tgbotapi.Update) {
	if upd.Message != nil {
		handleMessage(upd.Message)
	} else if upd.InlineQuery != nil {
		handleInlineQuery(upd.InlineQuery)
	}
}

func handleMessage(message *tgbotapi.Message) {
	code := message.Text

	if message.ReplyToMessage != nil {
		var nothing interface{}
		err := json.Unmarshal([]byte(message.ReplyToMessage.Text), &nothing)
		if err == nil {
			// means it's replying to a json message, so use that as input
			code = message.ReplyToMessage.Text + " | " + code
		}
	}

	ret, err := runjq(code)
	if err != nil {
		log.Warn().Err(err).Msg("message runjq")
		sendMessageAsReply(message.Chat.ID, `<pre><code>`+err.Error()+`</code></pre>`, message.MessageID)
		return
	}

	sendMessageAsReply(message.Chat.ID, `<pre><code class="language-json">`+ret+`</code></pre>`, message.MessageID)
}

func handleInlineQuery(q *tgbotapi.InlineQuery) {
	code := q.Query

	ret, err := runjq(code)
	if err != nil {
		log.Warn().Err(err).Msg("inline runjq")
		return
	}

	_, err = bot.AnswerInlineQuery(tgbotapi.InlineConfig{
		InlineQueryID: q.ID,
		Results: []interface{}{
			tgbotapi.NewInlineQueryResultArticleHTML("result", ret,
				`<pre><code class="language-json">`+ret+`</code></pre>`,
			),
			tgbotapi.NewInlineQueryResultArticleHTML("full", code+" → "+ret,
				`<code>`+code+`</code> => <code class="language-json">`+ret+`</code>`,
			),
		},
		IsPersonal: false,
		CacheTime:  30,
	})
	if err != nil {
		log.Warn().Err(err).Msg("inline results")
	}
}
