# README

## Introduction

Journal Agent is an innovative tool designed to streamline your note-taking and journaling experience. It automatically imports notes from AppleNotes, summarizes them using OpenAI's powerful AI technology, and then seamlessly integrates them into Notion. This tool is perfect for users who want to enhance their productivity and organize their notes efficiently.

## Features

- **Automatic Note Import**: Effortlessly imports all your notes from AppleNotes.
- **AI-Powered Summaries**: Utilizes OpenAI to create concise summaries of your notes.
- **Notion Integration**: Directly imports the summarized notes into your Notion pages.

## Configuration

Before running Journal Agent, you'll need to set up your configuration file. Here is the format for the configuration:

```yaml
note_dir: "path/to/your/notes"
api_key: "your_openai_api_key"
api_base: "openai_api_base_url"
notion_sk: "your_notion_secret_key"
page_id: "notion_page_id"
```

- `note_dir`: The directory where your AppleNotes are stored.
- `api_key`: Your OpenAI API key.
- `api_base`: The base URL for the OpenAI API.
- `notion_sk`: Your Notion integration secret key.
- `page_id`: The ID of the Notion page where you want to import the notes.

## How to Run

1. Clone the repository to your local machine.
2. Navigate to the cloned repository's directory.
3. Ensure that you have Go installed on your machine.
4. Create and set up the configuration file as described above.
5. Run the following command:

   ```shell
   go run main.go
   ```

## Prerequisites

- [Go](https://golang.org/dl/) programming language installed.
- Access to AppleNotes data.
- An OpenAI API key.
- A Notion account with API integration set up.

## Support

For support, please open an issue in the repository, and we will try to address it as soon as possible.

## Contributions

Contributions to Journal Agent are welcome! Please fork the repository and submit a pull request with your changes or improvements.

## License

Journal Agent is released under [MIT License](LICENSE.md).

---

Enjoy using Journal Agent to make your note management more efficient and productive!