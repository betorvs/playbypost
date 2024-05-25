async function decryptText(encryptedString, simpleKey) {
    // Convert the encryptedString to a byte array
    const encryptedBytes = new Uint8Array(atob(encryptedString).split("").map(c => c.charCodeAt(0)));
  
    // Convert the simple key to a UTF-8 encoded byte array
    const keyBytes = new TextEncoder().encode(simpleKey);
  
    // Create an AES-CBC cipher using the Web Crypto API
    const algorithm = {
      name: "AES-CBC",
      iv: new Uint8Array(16) // Replace with the IV used during encryption
    };
  
    // Import the key
    const key = await crypto.subtle.importKey(
      "raw",
      keyBytes,
      algorithm,
      false,
      ["decrypt"]
    );
  
    // Decrypt the encrypted bytes using the AES-CBC cipher and key
    const decryptedBytes = await crypto.subtle.decrypt(
      algorithm,
      key,
      encryptedBytes
    );
  
    // Convert the decrypted bytes to UTF-8 encoded string
    const decryptedString = new TextDecoder().decode(decryptedBytes);
  
    return decryptedString;
  }
  