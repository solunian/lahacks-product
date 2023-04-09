const RESPONSE_TYPE_NEW_TEXT = 0;
const RESPONSE_TYPE_FINAL = 1;
const AUTH_STRING = "LAH-SPRING-2023-9b2a97fa6ef69419"

// This is a client for a single message
class MessageClient {
  constructor(url, question, history) {
    this.url = url;
    this.question = question;
    this.history = history;
    this.source = null;
    this.eventListeners = {};
  }

  async connect() {
    const response = await fetch(this.url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': AUTH_STRING
      },
      body: JSON.stringify({
        question: this.question,
        history: this.history
      })
    });

    if (!response.ok) {
      throw new Error(`Server returned ${response.status}`);
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder('utf-8');
    let buffer = '';

    while (true) {
      const { done, value } = await reader.read();
      if (done) {
        break;
      }
      buffer += decoder.decode(value);
      const parts = buffer.split('\n');
      buffer = parts.pop();

      for (const part of parts) {
        if (part.length == 0) {
          continue;
        }

        const response = JSON.parse(part);
        switch (response.t) {
          case RESPONSE_TYPE_NEW_TEXT:
            this.trigger('newtext', response.n);
            break;
          case RESPONSE_TYPE_FINAL:
            this.trigger('finaltext', response.final_text, response.links);
            this.close();
            break;
          default:
            console.log('Unknown response type:', response.t);
        }
      }
    }
  }

  on(eventType, callback) {
    if (!this.eventListeners[eventType]) {
      this.eventListeners[eventType] = [];
    }
    this.eventListeners[eventType].push(callback);
  }

  off(eventType, callback) {
    if (this.eventListeners[eventType]) {
      const index = this.eventListeners[eventType].indexOf(callback);
      if (index !== -1) {
        this.eventListeners[eventType].splice(index, 1);
      }
    }
  }

  trigger(eventType, ...args) {
    if (this.eventListeners[eventType]) {
      for (const listener of this.eventListeners[eventType]) {
        listener.apply(this, args);
      }
    }
  }

  close() {
    if (this.source) {
      this.source.close();
    }
  }
}



// Example usage:
const client = new MessageClient('http://localhost:3000/conversation', 'The quick brown fox jumps over the lazy dog', [
  {
    human: "wifvfiu",
    assistant: "fwbfewjb"
  }
]);

client.connect();
client.on('newtext', (text) => console.log('Received new text:', text));
client.on('finaltext', (text, links) => {
    console.log('Received final:', text, links)
});
client.on('error', (event) => console.log('Error:', event));
client.on('close', () => console.log('Connection closed'));

// Here you might want to add an error when it times out
setTimeout(() => client.close(), 10000); // close the connection after 10 seconds