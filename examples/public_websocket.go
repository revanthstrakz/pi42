package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/net/websocket"
)

// ServerURL is the WebSocket server URL
const ServerURL = "wss://fawss.pi42.com/socket.io/?EIO=4&transport=websocket"

// Global variables
var (
	dbManager     *DatabaseManager
	dbInitialized bool
	conn          *websocket.Conn
	lastMessages  = make(map[string]json.RawMessage)
)

// DatabaseManager handles database operations
type DatabaseManager struct {
	db *sql.DB
}

// Symbol represents a trading symbol
type Symbol struct {
	Symbol            string `json:"symbol"`
	BaseAsset         string `json:"baseAsset"`
	QuoteAsset        string `json:"quoteAsset"`
	PricePrecision    int    `json:"pricePrecision"`
	QuantityPrecision int    `json:"quantityPrecision"`
	ContractName      string `json:"contractName"`
	ContractType      string `json:"contractType"`
	Market            string `json:"market"`
}

// ExchangeInfo represents the exchange information format
type ExchangeInfo struct {
	Markets   []string          `json:"markets"`
	Contracts []json.RawMessage `json:"contracts"`
}

// Initialize initializes the database connection
func (dm *DatabaseManager) Initialize() error {
	// In a real implementation, you would set up your database connection here
	// For this example, we'll just simulate success
	dm.db = &sql.DB{} // This is just a placeholder
	log.Println("Database initialized successfully")
	return nil
}

// SaveSymbols saves symbols to the database
func (dm *DatabaseManager) SaveSymbols(symbols []Symbol) error {
	// In a real implementation, you would save the symbols to the database
	// For this example, we'll just log the symbols
	log.Printf("Saving %d symbols to database", len(symbols))
	return nil
}

// Close closes the database connection
func (dm *DatabaseManager) Close() {
	// In a real implementation, you would close the database connection
	log.Println("Database connection closed")
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{}
}

// InitializeAndConnect initializes the database and connects to WebSocket
func InitializeAndConnect() {
	var err error

	// Initialize database
	log.Println("Initializing database...")
	err = dbManager.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	dbInitialized = true

	// Connect to WebSocket
	log.Println("Connecting to WebSocket server...")
	conn, err = connectWebSocket()
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}

	// Set up WebSocket handlers
	go handleWebSocketMessages()
}

// ConnectWebSocket connects to the WebSocket server
func connectWebSocket() (*websocket.Conn, error) {
	// Parse the URL
	u, err := url.Parse(ServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// Create WebSocket connection
	wsConfig, err := websocket.NewConfig(u.String(), "http://localhost/")
	if err != nil {
		return nil, fmt.Errorf("failed to create WebSocket config: %v", err)
	}

	conn, err := websocket.DialConfig(wsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial WebSocket: %v", err)
	}

	log.Println("Connected to WebSocket server")
	return conn, nil
}

// HandleWebSocketMessages handles incoming WebSocket messages
func handleWebSocketMessages() {
	// First, send the initial connection handshake
	err := websocket.Message.Send(conn, "40")
	if err != nil {
		log.Printf("Failed to send initial handshake: %v", err)
		return
	}

	// Now handle incoming messages
	for {
		var msg string
		err := websocket.Message.Receive(conn, &msg)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}

		// Process the message
		processWebSocketMessage(msg)
	}
}

// ProcessWebSocketMessage processes a WebSocket message
func processWebSocketMessage(msg string) {
	// Socket.IO protocol message parsing
	if len(msg) < 2 {
		return
	}

	// Handle different message types (simplified)
	switch {
	case msg == "0":
		// Socket.IO open packet
		log.Println("Socket.IO connection opened")
	case msg == "40":
		// Socket.IO connection established
		log.Println("Socket.IO connection established")

		// Fetch symbols and subscribe to topics
		symbols, err := fetchFuturesSymbols()
		if err != nil {
			log.Printf("Error fetching symbols: %v", err)
			return
		}

		if len(symbols) > 0 {
			// Save symbols to database
			err = dbManager.SaveSymbols(symbols)
			if err != nil {
				log.Printf("Error saving symbols to database: %v", err)
			}

			// Subscribe to topics
			subscribeToTopics(symbols)
		}
	case strings.HasPrefix(msg, "42"):
		// Socket.IO event
		eventData := msg[2:]
		var event []json.RawMessage
		err := json.Unmarshal([]byte(eventData), &event)
		if err != nil {
			log.Printf("Error parsing event: %v", err)
			return
		}

		if len(event) >= 2 {
			// Extract event name
			var eventName string
			err = json.Unmarshal(event[0], &eventName)
			if err != nil {
				log.Printf("Error parsing event name: %v", err)
				return
			}

			// Store the message data
			lastMessages[eventName] = event[1]
			log.Printf("%s event received", eventName)
		}
	}
}

// FetchFuturesSymbols fetches futures symbols from local JSON file
func fetchFuturesSymbols() ([]Symbol, error) {
	// Read exchange info from local file
	exchangeInfoPath := filepath.Join(".", "exchangeInfo.json")

	// Check if file exists
	_, err := os.Stat(exchangeInfoPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("exchangeInfo.json file does not exist")
	}

	// Read file
	fileData, err := os.ReadFile(exchangeInfoPath)
	if err != nil {
		return nil, fmt.Errorf("error reading exchangeInfo.json: %v", err)
	}

	if len(fileData) == 0 {
		return nil, fmt.Errorf("exchangeInfo.json file is empty")
	}

	log.Println("Reading exchange info from local file")

	// Parse JSON
	var exchangeInfo ExchangeInfo
	err = json.Unmarshal(fileData, &exchangeInfo)
	if err != nil {
		return nil, fmt.Errorf("error parsing exchangeInfo.json: %v", err)
	}

	if len(exchangeInfo.Contracts) == 0 {
		return nil, fmt.Errorf("no contracts found in exchangeInfo.json")
	}

	// Map contracts to symbols
	var symbols []Symbol
	for _, contractJSON := range exchangeInfo.Contracts {
		var contract map[string]interface{}
		err = json.Unmarshal(contractJSON, &contract)
		if err != nil {
			log.Printf("Error parsing contract: %v", err)
			continue
		}

		// Extract symbol information
		symbol := Symbol{
			Symbol:       contract["name"].(string),
			BaseAsset:    contract["baseAsset"].(string),
			QuoteAsset:   contract["quoteAsset"].(string),
			ContractName: contract["contractName"].(string),
		}

		// Parse precision values
		if pricePrecision, ok := contract["pricePrecision"].(string); ok {
			symbol.PricePrecision, _ = strconv.Atoi(pricePrecision)
		}

		if quantityPrecision, ok := contract["quantityPrecision"].(string); ok {
			symbol.QuantityPrecision, _ = strconv.Atoi(quantityPrecision)
		}

		// Add contract type and market if available
		if contractType, ok := contract["contractType"].(string); ok {
			symbol.ContractType = contractType
		}

		if market, ok := contract["market"].(string); ok {
			symbol.Market = market
		}

		symbols = append(symbols, symbol)
	}

	// Log sample of symbols
	log.Printf("Found %d symbols in local exchangeInfo.json:", len(symbols))
	sampleCount := 10
	if len(symbols) < sampleCount {
		sampleCount = len(symbols)
	}

	for i := 0; i < sampleCount; i++ {
		log.Printf("- %s (%s/%s)", symbols[i].Symbol, symbols[i].BaseAsset, symbols[i].QuoteAsset)
	}

	if len(symbols) > 10 {
		log.Printf("... and %d more symbols", len(symbols)-10)
	}

	return symbols, nil
}

// SubscribeToTopics subscribes to WebSocket topics for all symbols
func subscribeToTopics(symbols []Symbol) {
	if conn == nil {
		log.Println("WebSocket not initialized")
		return
	}

	// Define topics
	symbolTopics := []string{
		"depth_0.1",
		"markPrice",
		"kline_1m",
		"aggTrade",
		"ticker",
		"marketInfo",
	}

	allSymbolTopics := []string{
		"tickerArr",
		"markPriceArr",
		"allContractDetails",
	}

	// Create topics list
	var topics []string

	// Add symbol-specific topics
	for _, symbol := range symbols {
		lowerSymbol := strings.ToLower(symbol.Symbol)
		for _, topic := range symbolTopics {
			topics = append(topics, fmt.Sprintf("%s@%s", lowerSymbol, topic))
		}
	}

	// Add all-symbol topics
	topics = append(topics, allSymbolTopics...)

	log.Printf("Subscribing to %d topics", len(topics))

	// Create subscription message
	type SubscriptionParams struct {
		Params []string `json:"params"`
	}

	subscription := SubscriptionParams{
		Params: topics,
	}

	// Convert to JSON
	subscriptionJSON, err := json.Marshal(subscription)
	if err != nil {
		log.Printf("Error creating subscription JSON: %v", err)
		return
	}

	// Create Socket.IO event message
	eventMessage := fmt.Sprintf("42[\"subscribe\",%s]", string(subscriptionJSON))

	// Send subscription request
	err = websocket.Message.Send(conn, eventMessage)
	if err != nil {
		log.Printf("Error sending subscription request: %v", err)
		return
	}

	log.Println("Subscription request sent")
}

// SaveLastMessages saves the last messages to a file
func saveLastMessages() {
	filePath := filepath.Join(".", "lastNode.json")

	// Convert to JSON
	data, err := json.MarshalIndent(lastMessages, "", "  ")
	if err != nil {
		log.Printf("Error marshaling last messages: %v", err)
		return
	}

	// Write to file
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Printf("Error saving last messages: %v", err)
		return
	}

	log.Printf("Last messages saved to %s", filePath)
}

// PrintLastMessages prints the last messages to console
func printLastMessages() {
	log.Println("\n=== Last Received Messages ===")

	for eventType, message := range lastMessages {
		if len(message) > 0 {
			log.Printf("\n%s:", eventType)

			// Pretty-print JSON
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, message, "", "  ")
			if err != nil {
				log.Printf("Error formatting JSON: %v", err)
				log.Println(string(message))
			} else {
				log.Println(prettyJSON.String())
			}
		} else {
			log.Printf("\n%s: No message received", eventType)
		}
	}
}

func main() {
	log.Println("Starting futures data ingester...")

	// Create a database manager
	dbManager = NewDatabaseManager()

	// Initialize and connect
	InitializeAndConnect()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-sigChan
	log.Printf("\nReceived %s signal. Shutting down gracefully...", sig)

	// Print and save last messages
	printLastMessages()
	saveLastMessages()

	// Close database connection
	log.Println("Closing database connection...")
	dbManager.Close()

	// Close WebSocket connection
	if conn != nil {
		conn.Close()
	}

	log.Println("Exiting process")
}
