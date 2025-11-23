# Frontend Integration Guide - Create Vehicle Endpoint

## Endpoint: POST /api/user/vehicles

Complete guide for integrating the vehicle creation endpoint with your frontend application.

---

## Table of Contents

1. [Quick Start](#quick-start)
2. [Authentication Setup](#authentication-setup)
3. [Request Structure](#request-structure)
4. [Form Fields Reference](#form-fields-reference)
5. [Image Upload](#image-upload)
6. [Code Examples](#code-examples)
7. [Error Handling](#error-handling)
8. [Validation Rules](#validation-rules)
9. [Testing](#testing)

---

## Quick Start

### Endpoint Details

```
Method: POST
URL: http://localhost:8080/api/user/vehicles
Content-Type: multipart/form-data
Authorization: Bearer {JWT_TOKEN}
```

### Minimal Example

```javascript
const formData = new FormData();

// Required fields
formData.append('title', 'BMW 320d 2019');
formData.append('category', 'autoturisme');
formData.append('price', '18500');
formData.append('currency', 'lei');
formData.append('person_type', 'persoana_fizica');
formData.append('brand', 'BMW');
formData.append('model', '320d');
formData.append('fuel_type', 'motorina');
formData.append('body_type', 'sedan');
formData.append('year', '2019');
formData.append('condition', 'utilizat');
formData.append('transmission', 'automata');
formData.append('steering', 'stanga');
formData.append('city', 'București');
formData.append('contact_name', 'John Doe');
formData.append('email', 'john@example.com');

const response = await fetch('http://localhost:8080/api/user/vehicles', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
  },
  body: formData
});

const data = await response.json();
```

---

## Authentication Setup

### 1. Get JWT Token

Before creating a vehicle, user must be logged in:

```javascript
// Login to get token
const login = async (email, password) => {
  const response = await fetch('http://localhost:8080/api/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ email, password })
  });

  const data = await response.json();

  if (data.status === 'success') {
    // Store token in localStorage or state management
    localStorage.setItem('authToken', data.token);
    localStorage.setItem('userId', data.user.id);
    return data.token;
  }

  throw new Error(data.message || 'Login failed');
};
```

### 2. Check Authentication

```javascript
const isAuthenticated = () => {
  const token = localStorage.getItem('authToken');
  return !!token;
};

const getAuthToken = () => {
  return localStorage.getItem('authToken');
};
```

### 3. Handle Unauthorized Response

```javascript
if (response.status === 401) {
  // Token expired or invalid
  localStorage.removeItem('authToken');
  // Redirect to login page
  window.location.href = '/login';
}
```

---

## Request Structure

### Complete FormData Structure

```javascript
const createVehicleFormData = (vehicleData, images) => {
  const formData = new FormData();

  // Required fields
  formData.append('title', vehicleData.title);
  formData.append('category', vehicleData.category);
  formData.append('price', vehicleData.price);
  formData.append('currency', vehicleData.currency || 'lei');
  formData.append('person_type', vehicleData.personType);
  formData.append('brand', vehicleData.brand);
  formData.append('model', vehicleData.model);
  formData.append('fuel_type', vehicleData.fuelType);
  formData.append('body_type', vehicleData.bodyType);
  formData.append('year', vehicleData.year);
  formData.append('condition', vehicleData.condition);
  formData.append('transmission', vehicleData.transmission);
  formData.append('steering', vehicleData.steering);
  formData.append('city', vehicleData.city);
  formData.append('contact_name', vehicleData.contactName);
  formData.append('email', vehicleData.email);

  // Optional fields - only append if provided
  if (vehicleData.description) {
    formData.append('description', vehicleData.description);
  }

  if (vehicleData.negotiable !== undefined) {
    formData.append('negotiable', vehicleData.negotiable);
  }

  if (vehicleData.engineCapacity) {
    formData.append('engine_capacity', vehicleData.engineCapacity);
  }

  if (vehicleData.powerHp) {
    formData.append('power_hp', vehicleData.powerHp);
  }

  if (vehicleData.kilometers) {
    formData.append('kilometers', vehicleData.kilometers);
  }

  if (vehicleData.color) {
    formData.append('color', vehicleData.color);
  }

  if (vehicleData.numberOfKeys) {
    formData.append('number_of_keys', vehicleData.numberOfKeys);
  }

  if (vehicleData.registered !== undefined) {
    formData.append('registered', vehicleData.registered);
  }

  if (vehicleData.phone) {
    formData.append('phone', vehicleData.phone);
  }

  // Append images (max 8)
  if (images && images.length > 0) {
    images.forEach((image, index) => {
      if (index < 8) {
        formData.append('images', image);
      }
    });
  }

  return formData;
};
```

---

## Form Fields Reference

### Required Fields

| Field Name | Form Field | Type | Example | Notes |
|------------|------------|------|---------|-------|
| `title` | `title` | string | "BMW 320d 2019" | 5-255 characters |
| `category` | `category` | string | "autoturisme" | Vehicle category |
| `price` | `price` | number | 18500 | Must be > 0 |
| `currency` | `currency` | string | "lei" | Currency code |
| `person_type` | `person_type` | string | "persoana_fizica" | See options below |
| `brand` | `brand` | string | "BMW" | Brand name |
| `model` | `model` | string | "320d" | Model name |
| `fuel_type` | `fuel_type` | string | "motorina" | See options below |
| `body_type` | `body_type` | string | "sedan" | See options below |
| `year` | `year` | integer | 2019 | 1970-2030 |
| `condition` | `condition` | string | "utilizat" | "utilizat" or "nou" |
| `transmission` | `transmission` | string | "automata" | See options below |
| `steering` | `steering` | string | "stanga" | "stanga" or "dreapta" |
| `city` | `city` | string | "București" | City name |
| `contact_name` | `contact_name` | string | "John Doe" | Contact person |
| `email` | `email` | string | "john@example.com" | Valid email format |

### Optional Fields

| Field Name | Form Field | Type | Example | Notes |
|------------|------------|------|---------|-------|
| `description` | `description` | string | "Excellent condition..." | Long text |
| `negotiable` | `negotiable` | boolean | true | Price negotiable |
| `engine_capacity` | `engine_capacity` | integer | 1995 | In cm³ |
| `power_hp` | `power_hp` | integer | 190 | In HP |
| `kilometers` | `kilometers` | integer | 85000 | Mileage |
| `color` | `color` | string | "Black" | Color name |
| `number_of_keys` | `number_of_keys` | integer | 2 | Number of keys |
| `registered` | `registered` | boolean | true | Vehicle registered |
| `phone` | `phone` | string | "0721123456" | Phone number |
| `images` | `images` | file[] | File objects | Max 8 images |

### Field Options

#### Person Type (`person_type`)
```javascript
const personTypeOptions = [
  { value: 'persoana_fizica', label: 'Persoană Fizică' },
  { value: 'firma', label: 'Firmă' }
];
```

#### Fuel Type (`fuel_type`)
```javascript
const fuelTypeOptions = [
  { value: 'benzina', label: 'Benzină' },
  { value: 'motorina', label: 'Motorină' },
  { value: 'electric', label: 'Electric' },
  { value: 'hibrid', label: 'Hibrid' },
  { value: 'gpl', label: 'GPL' },
  { value: 'hybrid_benzina', label: 'Hybrid Benzină' },
  { value: 'hybrid_motorina', label: 'Hybrid Motorină' }
];
```

#### Body Type (`body_type`)
```javascript
const bodyTypeOptions = [
  { value: 'sedan', label: 'Sedan' },
  { value: 'suv', label: 'SUV' },
  { value: 'break', label: 'Break' },
  { value: 'coupe', label: 'Coupe' },
  { value: 'cabrio', label: 'Cabrio' },
  { value: 'hatchback', label: 'Hatchback' },
  { value: 'pickup', label: 'Pickup' },
  { value: 'van', label: 'Van' },
  { value: 'monovolum', label: 'Monovolum' }
];
```

#### Condition (`condition`)
```javascript
const conditionOptions = [
  { value: 'nou', label: 'Nou' },
  { value: 'utilizat', label: 'Utilizat' }
];
```

#### Transmission (`transmission`)
```javascript
const transmissionOptions = [
  { value: 'manuala', label: 'Manuală' },
  { value: 'automata', label: 'Automată' }
];
```

#### Steering (`steering`)
```javascript
const steeringOptions = [
  { value: 'stanga', label: 'Stânga' },
  { value: 'dreapta', label: 'Dreapta' }
];
```

---

## Image Upload

### Image Requirements

- **Format**: JPEG, PNG, JPG
- **Maximum**: 8 images per vehicle
- **Recommended size**: < 5MB per image
- **Field name**: `images` (multiple files with same name)

### React Example - File Input

```jsx
import { useState } from 'react';

const VehicleImageUpload = ({ onChange }) => {
  const [previews, setPreviews] = useState([]);

  const handleFileChange = (e) => {
    const files = Array.from(e.target.files);

    // Limit to 8 images
    const limitedFiles = files.slice(0, 8);

    // Create preview URLs
    const previewUrls = limitedFiles.map(file => URL.createObjectURL(file));
    setPreviews(previewUrls);

    // Pass files to parent
    onChange(limitedFiles);
  };

  return (
    <div>
      <input
        type="file"
        multiple
        accept="image/jpeg,image/png,image/jpg"
        onChange={handleFileChange}
        max="8"
      />

      <div className="previews">
        {previews.map((url, index) => (
          <img key={index} src={url} alt={`Preview ${index + 1}`} />
        ))}
      </div>

      <p>{previews.length} / 8 images selected</p>
    </div>
  );
};
```

### Vue 3 Example - File Input

```vue
<template>
  <div>
    <input
      type="file"
      multiple
      accept="image/jpeg,image/png,image/jpg"
      @change="handleFileChange"
      ref="fileInput"
    />

    <div class="previews">
      <img
        v-for="(url, index) in previews"
        :key="index"
        :src="url"
        :alt="`Preview ${index + 1}`"
      />
    </div>

    <p>{{ previews.length }} / 8 images selected</p>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const previews = ref([]);
const selectedFiles = ref([]);

const emit = defineEmits(['update:files']);

const handleFileChange = (event) => {
  const files = Array.from(event.target.files);

  // Limit to 8 images
  selectedFiles.value = files.slice(0, 8);

  // Create previews
  previews.value = selectedFiles.value.map(file =>
    URL.createObjectURL(file)
  );

  // Emit to parent
  emit('update:files', selectedFiles.value);
};
</script>
```

---

## Code Examples

### React + Axios

```jsx
import React, { useState } from 'react';
import axios from 'axios';

const CreateVehicleForm = () => {
  const [formData, setFormData] = useState({
    title: '',
    category: 'autoturisme',
    price: '',
    currency: 'lei',
    personType: 'persoana_fizica',
    brand: '',
    model: '',
    fuelType: 'benzina',
    bodyType: 'sedan',
    year: new Date().getFullYear(),
    condition: 'utilizat',
    transmission: 'automata',
    steering: 'stanga',
    city: '',
    contactName: '',
    email: '',
    // Optional fields
    description: '',
    negotiable: false,
    engineCapacity: '',
    powerHp: '',
    kilometers: '',
    color: '',
    numberOfKeys: '',
    registered: false,
    phone: ''
  });

  const [images, setImages] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const token = localStorage.getItem('authToken');

      if (!token) {
        throw new Error('Please login first');
      }

      const data = new FormData();

      // Append all form fields
      Object.keys(formData).forEach(key => {
        if (formData[key] !== '' && formData[key] !== null) {
          // Convert camelCase to snake_case for API
          const apiKey = key.replace(/([A-Z])/g, '_$1').toLowerCase();
          data.append(apiKey, formData[key]);
        }
      });

      // Append images
      images.forEach(image => {
        data.append('images', image);
      });

      const response = await axios.post(
        'http://localhost:8080/api/user/vehicles',
        data,
        {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'multipart/form-data'
          }
        }
      );

      console.log('Vehicle created:', response.data);

      // Redirect to vehicle page or show success message
      alert('Vehicle created successfully!');

      // Reset form
      setFormData({...initialFormData});
      setImages([]);

    } catch (err) {
      console.error('Error creating vehicle:', err);

      if (err.response?.status === 401) {
        setError('Please login to create a vehicle');
        // Redirect to login
      } else if (err.response?.status === 400) {
        setError(err.response.data.errors || err.response.data.message);
      } else {
        setError('Failed to create vehicle. Please try again.');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {error && <div className="error">{JSON.stringify(error)}</div>}

      {/* Form fields */}
      <input
        type="text"
        placeholder="Title"
        value={formData.title}
        onChange={(e) => setFormData({...formData, title: e.target.value})}
        required
      />

      {/* Add other fields... */}

      <input
        type="file"
        multiple
        accept="image/*"
        onChange={(e) => setImages(Array.from(e.target.files))}
      />

      <button type="submit" disabled={loading}>
        {loading ? 'Creating...' : 'Create Vehicle'}
      </button>
    </form>
  );
};

export default CreateVehicleForm;
```

### Vue 3 + Fetch

```vue
<template>
  <form @submit.prevent="handleSubmit">
    <div v-if="error" class="error">{{ error }}</div>

    <input
      v-model="formData.title"
      type="text"
      placeholder="Title"
      required
    />

    <input
      v-model="formData.brand"
      type="text"
      placeholder="Brand"
      required
    />

    <!-- Add other fields... -->

    <input
      type="file"
      multiple
      accept="image/*"
      @change="handleImageChange"
    />

    <button type="submit" :disabled="loading">
      {{ loading ? 'Creating...' : 'Create Vehicle' }}
    </button>
  </form>
</template>

<script setup>
import { ref } from 'vue';

const formData = ref({
  title: '',
  brand: '',
  model: '',
  price: '',
  // ... other fields
});

const images = ref([]);
const loading = ref(false);
const error = ref(null);

const handleImageChange = (event) => {
  images.value = Array.from(event.target.files);
};

const handleSubmit = async () => {
  loading.value = true;
  error.value = null;

  try {
    const token = localStorage.getItem('authToken');

    if (!token) {
      throw new Error('Please login first');
    }

    const data = new FormData();

    // Append form fields
    Object.keys(formData.value).forEach(key => {
      if (formData.value[key]) {
        const apiKey = key.replace(/([A-Z])/g, '_$1').toLowerCase();
        data.append(apiKey, formData.value[key]);
      }
    });

    // Append images
    images.value.forEach(image => {
      data.append('images', image);
    });

    const response = await fetch('http://localhost:8080/api/user/vehicles', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: data
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to create vehicle');
    }

    const result = await response.json();
    console.log('Vehicle created:', result);

    // Success - redirect or show message
    alert('Vehicle created successfully!');

  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
};
</script>
```

### Vanilla JavaScript

```javascript
document.getElementById('vehicleForm').addEventListener('submit', async (e) => {
  e.preventDefault();

  const token = localStorage.getItem('authToken');

  if (!token) {
    alert('Please login first');
    return;
  }

  const formData = new FormData();

  // Get form values
  formData.append('title', document.getElementById('title').value);
  formData.append('category', document.getElementById('category').value);
  formData.append('price', document.getElementById('price').value);
  formData.append('currency', 'lei');
  formData.append('person_type', document.getElementById('personType').value);
  formData.append('brand', document.getElementById('brand').value);
  formData.append('model', document.getElementById('model').value);
  formData.append('fuel_type', document.getElementById('fuelType').value);
  formData.append('body_type', document.getElementById('bodyType').value);
  formData.append('year', document.getElementById('year').value);
  formData.append('condition', document.getElementById('condition').value);
  formData.append('transmission', document.getElementById('transmission').value);
  formData.append('steering', document.getElementById('steering').value);
  formData.append('city', document.getElementById('city').value);
  formData.append('contact_name', document.getElementById('contactName').value);
  formData.append('email', document.getElementById('email').value);

  // Append images
  const imageFiles = document.getElementById('images').files;
  for (let i = 0; i < imageFiles.length && i < 8; i++) {
    formData.append('images', imageFiles[i]);
  }

  try {
    const response = await fetch('http://localhost:8080/api/user/vehicles', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: formData
    });

    const data = await response.json();

    if (response.ok) {
      alert('Vehicle created successfully!');
      console.log('Created vehicle:', data);
      // Redirect or reset form
    } else {
      alert('Error: ' + data.message);
    }
  } catch (error) {
    console.error('Error:', error);
    alert('Failed to create vehicle');
  }
});
```

---

## Error Handling

### Response Status Codes

| Status | Meaning | Action |
|--------|---------|--------|
| 201 | Created | Success! Vehicle created |
| 400 | Bad Request | Show validation errors |
| 401 | Unauthorized | Redirect to login |
| 500 | Server Error | Show generic error message |

### Error Response Format

```javascript
// 400 Bad Request - Validation Errors
{
  "status": "error",
  "message": "Validation failed",
  "errors": {
    "title": "Title must be between 5 and 255 characters",
    "price": "Price must be greater than 0",
    "email": "Invalid email format"
  }
}

// 401 Unauthorized
{
  "status": "error",
  "message": "User not authenticated"
}

// 500 Internal Server Error
{
  "status": "error",
  "message": "Failed to create vehicle",
  "error": "Database connection failed"
}
```

### Error Handling Example

```javascript
const handleApiError = (error, response) => {
  if (response?.status === 400 && response.data?.errors) {
    // Validation errors - show field-specific errors
    return {
      type: 'validation',
      errors: response.data.errors
    };
  }

  if (response?.status === 401) {
    // Unauthorized - redirect to login
    localStorage.removeItem('authToken');
    window.location.href = '/login';
    return {
      type: 'auth',
      message: 'Please login to continue'
    };
  }

  if (response?.status === 500) {
    // Server error
    return {
      type: 'server',
      message: 'Something went wrong. Please try again later.'
    };
  }

  // Generic error
  return {
    type: 'generic',
    message: error.message || 'An error occurred'
  };
};
```

---

## Validation Rules

### Client-Side Validation (Recommended)

```javascript
const validateVehicleForm = (formData) => {
  const errors = {};

  // Title
  if (!formData.title || formData.title.length < 5 || formData.title.length > 255) {
    errors.title = 'Title must be between 5 and 255 characters';
  }

  // Price
  if (!formData.price || formData.price <= 0) {
    errors.price = 'Price must be greater than 0';
  }

  // Year
  const currentYear = new Date().getFullYear();
  if (formData.year < 1970 || formData.year > currentYear + 1) {
    errors.year = `Year must be between 1970 and ${currentYear + 1}`;
  }

  // Email
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(formData.email)) {
    errors.email = 'Invalid email format';
  }

  // Required fields
  const requiredFields = [
    'title', 'category', 'price', 'currency', 'personType',
    'brand', 'model', 'fuelType', 'bodyType', 'year',
    'condition', 'transmission', 'steering', 'city',
    'contactName', 'email'
  ];

  requiredFields.forEach(field => {
    if (!formData[field]) {
      errors[field] = 'This field is required';
    }
  });

  return {
    isValid: Object.keys(errors).length === 0,
    errors
  };
};

// Usage
const { isValid, errors } = validateVehicleForm(formData);
if (!isValid) {
  setFormErrors(errors);
  return;
}
```

---

## Testing

### Test Data

```javascript
const testVehicle = {
  title: 'BMW 320d 2019 Test',
  category: 'autoturisme',
  description: 'This is a test vehicle',
  price: 18500,
  currency: 'lei',
  negotiable: true,
  personType: 'persoana_fizica',
  brand: 'BMW',
  model: '320d',
  engineCapacity: 1995,
  powerHp: 190,
  fuelType: 'motorina',
  bodyType: 'sedan',
  kilometers: 85000,
  color: 'Black',
  year: 2019,
  numberOfKeys: 2,
  condition: 'utilizat',
  transmission: 'automata',
  steering: 'stanga',
  registered: true,
  city: 'București',
  contactName: 'Test User',
  email: 'test@example.com',
  phone: '0721123456'
};
```

### cURL Test Command

```bash
curl -X POST http://localhost:8080/api/user/vehicles \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -F "title=BMW 320d 2019 Test" \
  -F "category=autoturisme" \
  -F "price=18500" \
  -F "currency=lei" \
  -F "person_type=persoana_fizica" \
  -F "brand=BMW" \
  -F "model=320d" \
  -F "fuel_type=motorina" \
  -F "body_type=sedan" \
  -F "year=2019" \
  -F "condition=utilizat" \
  -F "transmission=automata" \
  -F "steering=stanga" \
  -F "city=București" \
  -F "contact_name=Test User" \
  -F "email=test@example.com"
```

### Postman Collection

```json
{
  "info": {
    "name": "Create Vehicle",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Create Vehicle",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{authToken}}",
            "type": "text"
          }
        ],
        "body": {
          "mode": "formdata",
          "formdata": [
            {"key": "title", "value": "BMW 320d 2019", "type": "text"},
            {"key": "category", "value": "autoturisme", "type": "text"},
            {"key": "price", "value": "18500", "type": "text"},
            {"key": "currency", "value": "lei", "type": "text"},
            {"key": "person_type", "value": "persoana_fizica", "type": "text"},
            {"key": "brand", "value": "BMW", "type": "text"},
            {"key": "model", "value": "320d", "type": "text"},
            {"key": "fuel_type", "value": "motorina", "type": "text"},
            {"key": "body_type", "value": "sedan", "type": "text"},
            {"key": "year", "value": "2019", "type": "text"},
            {"key": "condition", "value": "utilizat", "type": "text"},
            {"key": "transmission", "value": "automata", "type": "text"},
            {"key": "steering", "value": "stanga", "type": "text"},
            {"key": "city", "value": "București", "type": "text"},
            {"key": "contact_name", "value": "Test User", "type": "text"},
            {"key": "email", "value": "test@example.com", "type": "text"}
          ]
        },
        "url": {
          "raw": "http://localhost:8080/api/user/vehicles",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "user", "vehicles"]
        }
      }
    }
  ]
}
```

---

## Success Response

```json
{
  "status": "success",
  "message": "Vehicle added successfully",
  "vehicle_id": 11,
  "data": {
    "id": 11,
    "user_id": 3,
    "uuid": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "slug": "bmw-320d-2019-test",
    "title": "BMW 320d 2019 Test",
    "category": "autoturisme",
    "description": "This is a test vehicle",
    "price": 18500.00,
    "currency": "lei",
    "negotiable": true,
    "person_type_id": 1,
    "person_type": "Persoană Fizică",
    "brand": "BMW",
    "model": "320d",
    "engine_capacity": 1995,
    "power_hp": 190,
    "fuel_type_id": 2,
    "fuel_type": "Motorină",
    "body_type_id": 1,
    "body_type": "Sedan",
    "kilometers": 85000,
    "color": "Black",
    "year": 2019,
    "number_of_keys": 2,
    "condition_id": 1,
    "condition": "Utilizat",
    "transmission_id": 2,
    "transmission": "Automată",
    "steering_id": 1,
    "steering": "Stânga",
    "registered": true,
    "city": "București",
    "contact_name": "Test User",
    "email": "test@example.com",
    "phone": "0721123456",
    "images": [
      {
        "id": 15,
        "vehicle_id": 11,
        "image_url": "/uploads/vehicles/uuid_timestamp_1.jpg",
        "created_at": "2025-11-23T10:30:00Z"
      }
    ],
    "created_at": "2025-11-23T10:30:00Z",
    "updated_at": "2025-11-23T10:30:00Z"
  }
}
```

---

## Common Issues & Solutions

### Issue: 401 Unauthorized

**Cause**: Missing or invalid token

**Solution**:
```javascript
// Check token exists
const token = localStorage.getItem('authToken');
if (!token) {
  // Redirect to login
  window.location.href = '/login';
}

// Include token in request
headers: {
  'Authorization': `Bearer ${token}`
}
```

### Issue: 400 Validation Error

**Cause**: Invalid form data

**Solution**: Validate form on client-side before sending

### Issue: Network Error

**Cause**: Backend not running or CORS issue

**Solution**:
- Check backend is running on port 8080
- Ensure CORS is configured for your frontend domain
- Check browser console for CORS errors

### Issue: Images not uploading

**Cause**: Wrong field name or file type

**Solution**:
```javascript
// Correct way to append multiple images
images.forEach(image => {
  formData.append('images', image); // Use 'images' (plural)
});

// Accept only jpeg, png, jpg
<input type="file" accept="image/jpeg,image/png,image/jpg" />
```

---

## Next Steps

1. **Implement Login Flow**: Set up authentication in your frontend
2. **Create Vehicle Form**: Build form with all required fields
3. **Add Image Upload**: Implement multi-image upload with preview
4. **Validation**: Add client-side validation
5. **Error Handling**: Display errors to user
6. **Success Handling**: Redirect after successful creation
7. **Testing**: Test with various data combinations

---

## Additional Resources

- **Swagger Documentation**: http://localhost:8080/swagger/index.html
- **API Endpoints Summary**: `/API_ENDPOINTS_SUMMARY.md`
- **Backend Repository**: Current project
- **Support**: GitHub Issues

---

**Last Updated**: 2025-11-23
**Endpoint**: POST /api/user/vehicles
**Version**: 1.0
