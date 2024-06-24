<?php
$servername = "viaduct.proxy.rlwy.net";
$username = "root";
$password = "MuLFjpjCHAKGLKBkOtXvIhWPbBIrdbAD";
$dbname = "railway";
$port = "41263"; // Add this line

// Create connection
$conn = new mysqli($servername, $username, $password, $dbname, $port); // Add $port parameter

// Check connection
if ($conn->connect_error) {
    die("Connection failed: " . $conn->connect_error);
}

if ($_SERVER["REQUEST_METHOD"] == "POST") {
    $name = $_POST['name'];
    $whatsapp = $_POST['whatsapp'];
    $message = $_POST['message'];

    $stmt = $conn->prepare("INSERT INTO submissions (name, whatsapp, message) VALUES (?, ?, ?)");
    $stmt->bind_param("sss", $name, $whatsapp, $message);

    if ($stmt->execute()) {
        $response = array("status" => "success", "message" => "Data submitted successfully!");
    } else {
        $response = array("status" => "error", "message" => "Failed to submit data.");
    }

    $stmt->close();
    $conn->close();

    header('Content-Type: application/json');
    echo json_encode($response);
} else {
    header('HTTP/1.1 405 Method Not Allowed');
    header('Allow: POST');
}
?>

