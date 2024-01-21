let serverUrl = 'http://localhost:9000';
let currentDirectory = '';

function makeRequest(method, endpoint, queryParams = {}) {
    const url = new URL(serverUrl + endpoint);
    Object.keys(queryParams).forEach(key => url.searchParams.append(key, queryParams[key]));

    return fetch(url, {
        method: method,
    })
	.then(response => {
            if (!response.ok) {
		console.error('Error:', response.statusText);
		return null;
            }
            // Check if the response has content before attempting to parse as JSON
            const contentLength = response.headers.get('Content-Length');
            if (contentLength === '0' || response.status === 204) {
		return null;
            }
            return response.json();
	})
	.catch(error => {
            console.error('Error:', error);
            return null;
	});
}

function displayDirectories(directories) {
    const directoriesDiv = document.getElementById('directories');
    directoriesDiv.innerHTML = '';

    if (directories != null) {
        directories.forEach(directory => {
            const button = document.createElement('button');
            button.textContent = `${directory.name}`;
            button.onclick = () => {
                currentDirectory = `${currentDirectory}/${directory.name}`;
                listDirectories();
            };
            directoriesDiv.appendChild(button);

            // Add space after each button
            const space = document.createTextNode(' ');
            directoriesDiv.appendChild(space);
        });
    }
}

function displayFiles(files) {
    const filesDiv = document.getElementById('files');
    filesDiv.innerHTML = '';
    if (files != null) {
        files.forEach(file => {
            const button = document.createElement('button');
            button.textContent = `${file.name}`;
            filesDiv.appendChild(button);

            // Add space after each button
            const space = document.createTextNode(' ');
            filesDiv.appendChild(space);

            // Add right-click context menu
            button.addEventListener('contextmenu', (event) => {
                event.preventDefault();
                showContextMenu(event, file, filesDiv, button);
            });
        });
    }
}

function showContextMenu(event, file, filesDiv, button) {
    const contextMenu = document.createElement('div');
    contextMenu.className = 'context-menu';
    contextMenu.innerHTML = `
        <div class="context-menu-item" onclick="downloadFile('${file.name}')">Download</div>
        <div class="context-menu-item" onclick="showDetails('${file.name}')">Details</div>
        <div class="context-menu-item" onclick="deleteFile('${file.name}')">Delete</div>
    `;

    // Calculate the top position relative to the filesDiv
    const filesDivRect = filesDiv.getBoundingClientRect();
    const buttonRect = button.getBoundingClientRect();
    const contextMenuHeight = contextMenu.offsetHeight;
    const topPosition = buttonRect.top - filesDivRect.top - contextMenuHeight; // Remove the padding

    contextMenu.style.top = `${Math.max(topPosition, 0)}px`;
    contextMenu.style.left = `0`;

    filesDiv.appendChild(contextMenu);

    // Remove the context menu when clicking outside of it
    const removeContextMenu = () => {
	filesDiv.removeChild(contextMenu);
        if (filesDiv.contains(contextMenu)) {
            filesDiv.removeChild(contextMenu);
        }
        document.removeEventListener('click', removeContextMenu);
    };

    document.addEventListener('click', removeContextMenu);
}

function downloadFile(filename) {
    // Logic to download the file
    console.log(`Downloading file: ${filename}`);
}

function showDetails(filename) {
    // Logic to show details of the file
    console.log(`Details for file: ${filename}`);
}

function deleteFile(filename) {
    // Logic to delete the file
    console.log(`Deleting file: ${filename}`);
}

function listDirectories(path) {
    if (path != undefined)
	currentDirectory = path;
    console.log('Listing directory: ' + currentDirectory)
    makeRequest('GET', `/${currentDirectory}`, { type: 'directory', operation: 'list'})
        .then(data => {
	    var directories
	    var files
	    if ( data != null) {
		directories = data.filter(item => item.type === "directory");
		files = data.filter(item => item.type === "file");
	    }
	    displayDirectories(directories);
	    displayFiles(files);
	    displayBreadcrumbs();
        })
        .catch(error => console.error('Error:', error));
}

function displayBreadcrumbs() {
    const breadcrumbsDiv = document.getElementById('breadcrumbs');

    if (breadcrumbsDiv) {
        breadcrumbsDiv.innerHTML = '';

        const directories = currentDirectory.split('/');
        let path = '';

        const rootButton = document.createElement('button');
        rootButton.textContent = 'Drive';
        rootButton.onclick = () => {
            currentDirectory = '';
            listDirectories();
        };
        breadcrumbsDiv.appendChild(rootButton);

        for (let i = 0; i < directories.length; i++) {
            const directory = directories[i];
            if (directory) {
                path += `/${directory}`;
                const button = document.createElement('button');
                button.textContent = directory;
                const newPath = path;
                button.onclick = () => {
                    currentDirectory = newPath;
                    listDirectories();
                };

		const space = document.createTextNode(' > ');
		breadcrumbsDiv.appendChild(space);
                breadcrumbsDiv.appendChild(button);
            }
        }
    }
}


function listRecent() {
    // Logic for listing recent files or directories
    // Implement as per your API's functionality
}

function listStarred() {
    // Logic for listing starred files or directories
    // Implement as per your API's functionality
}

function newDirectory() {
    const directoryName = prompt('Enter new directory name:');
    if (directoryName) {
        makeRequest('POST', '/' + currentDirectory + '/' + directoryName, { type: 'directory'});
	listDirectories()
    }
}

function newFile() {
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.addEventListener('change', (event) => {
        const file = event.target.files[0];
        const filename = event.target.value.split(/(\\|\/)/g).pop(); // Extracts the filename from the path
        if (file) {
            uploadFile(file, filename);
        }
    });
    // Append the file input to the body to trigger the file selection dialog
    document.body.appendChild(fileInput);
    fileInput.click();
    // Remove the file input from the body after selection (optional)
    fileInput.addEventListener('change', () => {
        document.body.removeChild(fileInput);
    });

}

function uploadFile(file, filename) {
    console.log("Uploading -> "+filename)
    const formData = new FormData();
    formData.append('file', file);
    console.log(currentDirectory);

    queryParams = { type: 'file'}
    const url = new URL(serverUrl + '/' + currentDirectory + '/' + filename);
    Object.keys(queryParams).forEach(key => url.searchParams.append(key, queryParams[key]));

    fetch(url, {
        method: 'POST',
        body: formData,
        headers: {
            'Content-Type': 'application/octet-stream',
	    'Content-Disposition': `attachment; filename="${filename}"`
        }
    })
	.then(response => {
            if (response.ok) {
		// Handle successful response
		console.log('File uploaded successfully!');
		listDirectories()
            } else {
		// Handle error response
		console.error('Failed to upload file.');
            }
	})
	.catch(error => {
            // Handle fetch error
            console.error('Error:', error);
	});
}

function goToParent() {
    if (currentDirectory) {
        const lastIndex = currentDirectory.lastIndexOf('/');
        if (lastIndex !== -1) {
            currentDirectory = currentDirectory.substring(0, lastIndex);
            listDirectories();
        }
    }
}

// Initial call to set up the UI
listDirectories();
