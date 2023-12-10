package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/YuanData/allegro-trade/util"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Test Mail"
	content := "Test Mail"
	to := []string{"test@mail.tst"}
	
	err = sender.SendEmail(subject, content, to, nil, nil, nil)
	require.NoError(t, err)
}
