<h1>Push file (GitHub Repo Manager)</h1>

<p><strong>Push file (GitHub Repo Manager)</strong> is a command-line tool written in Go-Lang that makes creating GitHub repositories and making commits easier with just a single click. The tool is cross-platform and can be compiled for Windows, Linux, and macOS. With just one click, you can:</p>

<ul>
    <li>Create a new GitHub repository</li>
    <li>Initialize a local repository</li>
    <li>Add and commit your changes</li>
    <li>Push the commit to GitHub</li>
</ul>

<h2>Features</h2>

<ul>
    <li><strong>Cross-Platform</strong>: Compatible with Windows, Linux, and macOS.</li>
    <li><strong>Single-Click Actions</strong>: Perform GitHub repo management tasks in a streamlined, simplified manner.</li>
    <li><strong>Easy Compilation</strong>: Compile the program for any platform using pre-configured batch scripts.</li>
</ul>

<h2>Requirements</h2>

<ul>
    <li><a href="https://golang.org/dl/">Go</a> 1.16 or later</li>
    <li>GitHub account and <a href="https://cli.github.com/">GitHub CLI</a> installed and configured</li>
    <li>Git installed on your machine</li>
</ul>

<h2>Installation</h2>

<ol>
    <strong>Clone the repository</strong>:
        <pre><code>git clone https://github.com/ArminKardan/push.git
cd push
        </code></pre>
    
</ol>

<h2>Compilation</h2>

<h3>For Windows</h3>

<p>Simply run the <code>compile.bat</code> file:</p>

<pre><code>compile.bat</code></pre>

<p>This will compile the program for Windows, and the executable will be created in the <code>bin/windows/</code> folder.</p>

<h3>For All Platforms (Windows, Linux, macOS)</h3>

<p>To compile the program for all three platforms (Windows, Linux, and macOS), use the <code>compile-all.bat</code> script:</p>

<pre><code>compile-all.bat</code></pre>

<p>This will create executables for Windows, Linux, and macOS in the <code>bin/windows/</code>, <code>bin/linux/</code>, and <code>bin/mac/</code> folders, respectively.</p>

<h2>Usage</h2>

<p>Once compiled, you can run the executable from the command line or terminal.</p>

<ol>
    <li><strong>Create a New GitHub Repository and Make Initial Commit</strong>:
        <pre>Just rename project directory to your desire name, put the file inside the directory and run it!</pre>
        <p>This will:</p>
        <ul>
            <li>Create a new repository on your GitHub account with the specified <code>&lt;directory-name&gt;</code></li>
            <li>Initialize a local Git repository</li>
            <li>Add all files in the current directory</li>
            <li>Make the initial commit</li>
            <li>Push the changes to GitHub</li>
        </ul>
    </li>
</ol>

<h2>Contributing</h2>

<p>Contributions are welcome! Feel free to submit a pull request or open an issue if you encounter any bugs or have suggestions for improvement.</p>

<h2>License</h2>

<p>This project is licensed under the MIT License - see the <a href="LICENSE">LICENSE</a> file for details.</p>

<h2>Author</h2>

<ul>
    <li><strong>Ethan Cardan</strong> - <a href="https://github.com/arminkardan">GitHub Profile</a></li>
</ul>
