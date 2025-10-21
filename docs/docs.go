package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "email": "support@example.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/upload": {
            "post": {
                "description": "Upload a file to the object storage server. Supports images (with automatic resizing), videos (with transcoding), and audio files (with bitrate conversion). Maximum file size: 4GB per request.",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Upload a file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "File to upload",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "File uploaded successfully",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "success": {
                                    "type": "boolean",
                                    "example": true
                                },
                                "message": {
                                    "type": "string",
                                    "example": "File uploaded successfully"
                                },
                                "data": {
                                    "type": "object",
                                    "properties": {
                                        "file_id": {
                                            "type": "string",
                                            "example": "01932ed7-8b5c-7890-abcd-1234567890ab"
                                        },
                                        "filename": {
                                            "type": "string",
                                            "example": "example.jpg"
                                        },
                                        "size": {
                                            "type": "integer",
                                            "example": 1048576
                                        },
                                        "mime_type": {
                                            "type": "string",
                                            "example": "image/jpeg"
                                        },
                                        "urls": {
                                            "type": "object",
                                            "properties": {
                                                "metadata": {
                                                    "type": "string",
                                                    "example": "http://localhost:8080/api/files/01932ed7-8b5c-7890-abcd-1234567890ab"
                                                },
                                                "download": {
                                                    "type": "string",
                                                    "example": "http://localhost:8080/api/files/01932ed7-8b5c-7890-abcd-1234567890ab/download"
                                                },
                                                "view": {
                                                    "type": "object",
                                                    "properties": {
                                                        "original": {
                                                            "type": "string",
                                                            "example": "http://localhost:8080/api/files/01932ed7-8b5c-7890-abcd-1234567890ab/view"
                                                        },
                                                        "thumbnail": {
                                                            "type": "string",
                                                            "example": "http://localhost:8080/api/files/01932ed7-8b5c-7890-abcd-1234567890ab/view?resolution=thumbnail"
                                                        },
                                                        "small": {
                                                            "type": "string",
                                                            "example": "http://localhost:8080/api/files/01932ed7-8b5c-7890-abcd-1234567890ab/view?resolution=small"
                                                        },
                                                        "medium": {
                                                            "type": "string",
                                                            "example": "http://localhost:8080/api/files/01932ed7-8b5c-7890-abcd-1234567890ab/view?resolution=medium"
                                                        },
                                                        "large": {
                                                            "type": "string",
                                                            "example": "http://localhost:8080/api/files/01932ed7-8b5c-7890-abcd-1234567890ab/view?resolution=large"
                                                        }
                                                    }
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad request - no file uploaded or invalid file",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "success": {
                                    "type": "boolean",
                                    "example": false
                                },
                                "message": {
                                    "type": "string",
                                    "example": "No file uploaded"
                                }
                            }
                        }
                    },
                    "413": {
                        "description": "File too large (max 4GB)",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "success": {
                                    "type": "boolean",
                                    "example": false
                                },
                                "message": {
                                    "type": "string",
                                    "example": "File size exceeds limit"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "success": {
                                    "type": "boolean",
                                    "example": false
                                },
                                "message": {
                                    "type": "string",
                                    "example": "Failed to save file"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/files/{id}": {
            "get": {
                "description": "Get file metadata including URLs for download and view with different resolutions",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Get file metadata",
                "parameters": [
                    {
                        "type": "string",
                        "description": "File ID (UUID v7)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "File metadata",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "success": {
                                    "type": "boolean",
                                    "example": true
                                },
                                "data": {
                                    "type": "object",
                                    "properties": {
                                        "file_id": {
                                            "type": "string",
                                            "example": "01932ed7-8b5c-7890-abcd-1234567890ab"
                                        },
                                        "filename": {
                                            "type": "string",
                                            "example": "example.jpg"
                                        },
                                        "size": {
                                            "type": "integer",
                                            "example": 1048576
                                        },
                                        "mime_type": {
                                            "type": "string",
                                            "example": "image/jpeg"
                                        },
                                        "urls": {
                                            "type": "object"
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "404": {
                        "description": "File not found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "success": {
                                    "type": "boolean",
                                    "example": false
                                },
                                "message": {
                                    "type": "string",
                                    "example": "File not found"
                                }
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a file and all its variants (thumbnails, resolutions, etc)",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Delete a file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "File ID (UUID v7)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "File deleted successfully",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "success": {
                                    "type": "boolean",
                                    "example": true
                                },
                                "message": {
                                    "type": "string",
                                    "example": "File deleted successfully"
                                }
                            }
                        }
                    },
                    "404": {
                        "description": "File not found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "success": {
                                    "type": "boolean",
                                    "example": false
                                },
                                "message": {
                                    "type": "string",
                                    "example": "File not found"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/files/{id}/download": {
            "get": {
                "description": "Download the original file",
                "produces": [
                    "application/octet-stream"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Download a file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "File ID (UUID v7)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "File content",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "404": {
                        "description": "File not found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "success": {
                                    "type": "boolean",
                                    "example": false
                                },
                                "message": {
                                    "type": "string",
                                    "example": "File not found"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/files/{id}/view": {
            "get": {
                "description": "View a file with optional resolution parameter. For images: thumbnail, small, medium, large. For videos: 360p, 480p, 720p, 1080p. For audio: 64k, 128k, 320k.",
                "produces": [
                    "image/jpeg",
                    "image/png",
                    "video/mp4",
                    "audio/mpeg"
                ],
                "tags": [
                    "files"
                ],
                "summary": "View a file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "File ID (UUID v7)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Resolution (thumbnail/small/medium/large for images, 360p/480p/720p/1080p for videos, 64k/128k/320k for audio)",
                        "name": "resolution",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "File content",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "404": {
                        "description": "File not found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "success": {
                                    "type": "boolean",
                                    "example": false
                                },
                                "message": {
                                    "type": "string",
                                    "example": "File not found"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Check if the server is running",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "Server is healthy",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "status": {
                                    "type": "string",
                                    "example": "ok"
                                },
                                "timestamp": {
                                    "type": "string",
                                    "example": "2025-10-21T10:30:00Z"
                                }
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {}
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Object Storage API",
	Description:      "High-performance object storage server with automatic image resizing, video transcoding, and audio conversion. Supports up to 4GB files per request with streaming upload/download.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
