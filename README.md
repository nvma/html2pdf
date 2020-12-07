# html2pdf
Small utility for rapid prototyping of PDF layout based on HTML documents. Watcher will look for changes in the provided HTML file and will rebuild PDF on every change. 

## Usage
Linux: `./bin/html2pdf-linux -filename="template"`

MacOS: `./bin/html2pdf-darwin -filename="template"`

Win10: `./bin/html2pdf.exe -filename="template"`

`filename=""` flag accepts any HTML document located in the `html` directory, PDF output will be generated automatically with the same name and placed in the `pdf` directory.
