#include <Arduino.h>
#include <PubSubClient.h>
#include <WiFiClient.h>
#include <WiFiManager.h>

#define DEFAULT_PLACE_ID "111121314151617"

#define RELAY 0
#define RELAY_DELAY 1000

const char *mqttServer = "aaa.bbb.ccc.ddd";
int mqttPort = 1883;
const char *mqttUsername = "mqtt";
const char *mqttPassword = "mqtt";
String mqttTopic = "door/";

WiFiManagerParameter placeId("placeId", "Place ID", DEFAULT_PLACE_ID, 15);

WiFiClient wifiClient;
PubSubClient client(wifiClient);

long lastReconnectAttempt = 0;
long relayAttempt = 0;

void callback(char *topic, byte *payload, unsigned int length) {
  if (relayAttempt == 0) {
    digitalWrite(RELAY, HIGH);
    relayAttempt = millis();
  }
}

boolean reconnect() {
  String pId = String(placeId.getValue());
  String id =
      "ESP-" + String(placeId.getValue()) + "-" + String(ESP.getChipId(), HEX);
  String topic = mqttTopic + pId;
  if (client.connect(id.c_str(), mqttUsername, mqttPassword)) {

    client.subscribe(topic.c_str());
  }
  return client.connected();
}

void setup() {

  Serial.begin(115200);
  pinMode(RELAY, OUTPUT);

  WiFiManager wifiManager;
  wifiManager.addParameter(&placeId);
  wifiManager.autoConnect();

  if (strlen(placeId.getValue()) != 15) {
    Serial.println("Invalid place id\nResetting in 5 secs ...");
    delay(5000);
    wifiManager.resetSettings();
    ESP.restart();
  }

  client.setServer(mqttServer, mqttPort);
  client.setCallback(callback);

  lastReconnectAttempt = 0;
}

void loop() {
  if (!client.connected()) {
    long now = millis();
    if (now - lastReconnectAttempt > 5000) {
      lastReconnectAttempt = now;
      // Attempt to reconnect
      if (reconnect()) {
        lastReconnectAttempt = 0;
      }
    }
  } else {  
    client.loop();
  }

  if((relayAttempt != 0) && (millis() - relayAttempt) > RELAY_DELAY){
    digitalWrite(RELAY, LOW);
    relayAttempt = 0;
  }

}
