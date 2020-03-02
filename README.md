# goprogramming
To gain understanding of the GO language. And application of same in developing a very basic level http server.


### OCR
Extracts text from a .png file and prints it to console<br/>
lib - github.com/otiai10/gosseract

### OCR_Vision
Demonstrates the use of Google's Vision lib for OCR in a pdf. Extracts text from a pdf saved on Google Cloud.<br/>
lib - cloud.google.com/go/vision/apiv1

### PDF_Reader
Extracts text from a .pdf file and prints it to console.<br/>
lib - github.com/ledongthuc/pdf

### PDF_Images_Reader
List down the images in a pdf and print the respective data on console.<br/>
lib - github.com/unidoc/unidoc/pdf

### unipdf_reader
Extract images from pdf, prints the text to console and saves the extracted images as .png.<br/>
Arguments - 1. Input path to the pdf file. 2. Output path where the extracted images need to be stored.<br/>
lib - github.com/unidoc/unipdf/extractor

### unipdf_reader_with_download
Extract images from an online pdf after downloading it to local storage. Prints the extracted text to console and saves the images as .png.<br/>
Arguments - 1. URL of the online pdf. 2. Output path where the extracted images need to be stored.<br/>
lib - github.com/unidoc/unipdf/extractor

### http_server 
http://localhost:8000/fetch?date=03DEC2019<br/>
Demostrates the usage of net/http lib to develop a server.<br/>
Uses query param 'Date' to formulate a dynamic url that will download a zip folder containing a .csv file. The .csv file is further saved to local storage, content records are read and passes back as response.<br/>
The 'Date' param is validated for correct date format. In case of error, an error string is returned as a response.

