let serverUrl = 'http://localhost:9000';
let currentDirectory = '';

function setServerUrl() {
    serverUrl = document.getElementById('serverUrl').value;
    alert('Server URL set: ' + serverUrl);
    listDirectories();
}

function makeRequest(method, endpoint, queryParams = {}) {
    const url = new URL(serverUrl + endpoint);
    Object.keys(queryParams).forEach(key => url.searchParams.append(key, queryParams[key]));

    return fetch(url, {
        method: method,
    })
    .then(response => response.json())
    .catch(error => {
        console.error('Error:', error);
        return null;
    });
}

function displayDirectories(directories) {
    const directoriesDiv = document.getElementById('directories');
    directoriesDiv.innerHTML = '';

    if (directories !== null) {
        directories.forEach(directory => {
            const button = document.createElement('button');
            button.textContent = `📁 ${directory}`;
            button.onclick = () => {
                currentDirectory = `${currentDirectory}/${directory}`;
                listDirectories();
            };
            directoriesDiv.appendChild(button);
        });
    }
}

function listDirectories(path) {
    console.log('-> Listing directories ' + path)
    if (path !== undefined)
	currentDirectory = path;
    console.log('Listing directories ' + currentDirectory)
    makeRequest('GET', `/${currentDirectory}`, { type: 'directory', operation: 'list'})
        .then(data => {
            displayDirectories(data);
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
        rootButton.textContent = '/';
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
                breadcrumbsDiv.appendChild(button);
            }
        }
    }
}


function getFiles() {
    makeRequest('GET', '/GetFile', { type: 'file' });
}

function listRecent() {
    // Logic for listing recent files or directories
    // Implement as per your API's functionality
}

function listStarred() {
    // Logic for listing starred files or directories
    // Implement as per your API's functionality
}

function newOption() {
    const option = prompt('Enter "directory" to create a new directory or "file" to upload a new file:');
    if (option === 'directory') {
        newDirectory();
    } else if (option === 'file') {
        newFile();
    } else {
        alert('Invalid option entered. Please enter "directory" or "file".');
    }
}

function newDirectory() {
    const directoryName = prompt('Enter new directory name:');
    if (directoryName) {
        makeRequest('POST', '/CreateDirectory', { type: 'directory', name: directoryName });
    }
}

function newFile() {
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.addEventListener('change', (event) => {
        const file = event.target.files[0];
        if (file) {
            uploadFile(file);
        }
    });
    fileInput.click();
}

function uploadFile(file) {
    const formData = new FormData();
    formData.append('file', file);
    console.log(currentDirectory);

    return fetch(serverUrl + '/' + currentDirectory + '/' + file.name + '?type=file', {
        method: 'POST',
        body: formData,
    })
    .then(response => response.json())
    .then(data => {
        console.log(data);
        // Handle the response data here
        // You can update the UI or perform further actions based on the response
    })
    .catch(error => console.error('Error:', error));
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
