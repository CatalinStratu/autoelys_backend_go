# Quick Reference - Create Vehicle Endpoint

## ⚡ Quick Copy-Paste Code

### Minimal Working Example

```javascript
// 1. Get token from login
const token = localStorage.getItem('authToken');

// 2. Create FormData
const formData = new FormData();
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

// 3. Send request
const response = await fetch('http://localhost:8080/api/user/vehicles', {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` },
  body: formData
});

const data = await response.json();
console.log(data);
```

---

## 📋 Required Fields Checklist

```javascript
const requiredFields = {
  title: 'Vehicle title (5-255 chars)',
  category: 'autoturisme',
  price: 'Number > 0',
  currency: 'lei',
  person_type: 'persoana_fizica | firma',
  brand: 'Brand name',
  model: 'Model name',
  fuel_type: 'benzina | motorina | electric | hibrid | gpl',
  body_type: 'sedan | suv | break | coupe | cabrio | hatchback | pickup | van | monovolum',
  year: '1970-2030',
  condition: 'utilizat | nou',
  transmission: 'manuala | automata',
  steering: 'stanga | dreapta',
  city: 'City name',
  contact_name: 'Contact person',
  email: 'valid@email.com'
};
```

---

## 🎨 Dropdown Options

```javascript
// Copy these for your select inputs
const dropdowns = {
  personType: [
    { value: 'persoana_fizica', label: 'Persoană Fizică' },
    { value: 'firma', label: 'Firmă' }
  ],

  fuelType: [
    { value: 'benzina', label: 'Benzină' },
    { value: 'motorina', label: 'Motorină' },
    { value: 'electric', label: 'Electric' },
    { value: 'hibrid', label: 'Hibrid' },
    { value: 'gpl', label: 'GPL' }
  ],

  bodyType: [
    { value: 'sedan', label: 'Sedan' },
    { value: 'suv', label: 'SUV' },
    { value: 'break', label: 'Break' },
    { value: 'coupe', label: 'Coupe' },
    { value: 'cabrio', label: 'Cabrio' },
    { value: 'hatchback', label: 'Hatchback' },
    { value: 'pickup', label: 'Pickup' },
    { value: 'van', label: 'Van' },
    { value: 'monovolum', label: 'Monovolum' }
  ],

  condition: [
    { value: 'nou', label: 'Nou' },
    { value: 'utilizat', label: 'Utilizat' }
  ],

  transmission: [
    { value: 'manuala', label: 'Manuală' },
    { value: 'automata', label: 'Automată' }
  ],

  steering: [
    { value: 'stanga', label: 'Stânga' },
    { value: 'dreapta', label: 'Dreapta' }
  ]
};
```

---

## 🖼️ Image Upload

```javascript
// HTML
<input
  type="file"
  multiple
  accept="image/jpeg,image/png,image/jpg"
  onChange={handleImageChange}
  max="8"
/>

// JavaScript
const handleImageChange = (e) => {
  const files = Array.from(e.target.files);
  const limitedFiles = files.slice(0, 8); // Max 8 images

  // Add to FormData
  limitedFiles.forEach(file => {
    formData.append('images', file);
  });
};
```

---

## ✅ Validation

```javascript
const validate = (data) => {
  const errors = {};

  if (!data.title || data.title.length < 5 || data.title.length > 255) {
    errors.title = 'Title: 5-255 characters';
  }

  if (!data.price || data.price <= 0) {
    errors.price = 'Price must be > 0';
  }

  if (data.year < 1970 || data.year > new Date().getFullYear() + 1) {
    errors.year = 'Invalid year';
  }

  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email)) {
    errors.email = 'Invalid email';
  }

  return Object.keys(errors).length === 0 ? null : errors;
};
```

---

## 🚨 Error Handling

```javascript
try {
  const response = await fetch(url, { method: 'POST', headers, body });
  const data = await response.json();

  if (!response.ok) {
    if (response.status === 401) {
      // Redirect to login
      window.location.href = '/login';
    } else if (response.status === 400) {
      // Show validation errors
      console.error('Validation errors:', data.errors);
    } else {
      // Generic error
      console.error('Error:', data.message);
    }
  } else {
    // Success!
    console.log('Vehicle created:', data.data);
  }
} catch (error) {
  console.error('Network error:', error);
}
```

---

## 📦 Response Format

### Success (201)
```json
{
  "status": "success",
  "message": "Vehicle added successfully",
  "vehicle_id": 11,
  "data": {
    "id": 11,
    "user_id": 3,
    "uuid": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "slug": "bmw-320d-2019",
    "title": "BMW 320d 2019",
    "price": 18500.00,
    ...
  }
}
```

### Error (400)
```json
{
  "status": "error",
  "message": "Validation failed",
  "errors": {
    "title": "Title must be between 5 and 255 characters",
    "price": "Price must be greater than 0"
  }
}
```

---

## 🔑 Authentication

```javascript
// Login first
const login = async (email, password) => {
  const res = await fetch('http://localhost:8080/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  const data = await res.json();
  localStorage.setItem('authToken', data.token);
  return data.token;
};

// Use token
const token = localStorage.getItem('authToken');
headers: { 'Authorization': `Bearer ${token}` }
```

---

## 📝 Complete React Example

```jsx
import { useState } from 'react';

const CreateVehicle = () => {
  const [form, setForm] = useState({
    title: '', brand: '', model: '', price: '',
    year: 2024, city: '', contactName: '', email: ''
  });
  const [images, setImages] = useState([]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem('authToken');
    const data = new FormData();

    // Required fields
    data.append('title', form.title);
    data.append('category', 'autoturisme');
    data.append('price', form.price);
    data.append('currency', 'lei');
    data.append('person_type', 'persoana_fizica');
    data.append('brand', form.brand);
    data.append('model', form.model);
    data.append('fuel_type', 'benzina');
    data.append('body_type', 'sedan');
    data.append('year', form.year);
    data.append('condition', 'utilizat');
    data.append('transmission', 'automata');
    data.append('steering', 'stanga');
    data.append('city', form.city);
    data.append('contact_name', form.contactName);
    data.append('email', form.email);

    images.forEach(img => data.append('images', img));

    const res = await fetch('http://localhost:8080/api/user/vehicles', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` },
      body: data
    });

    const result = await res.json();
    console.log(result);
  };

  return (
    <form onSubmit={handleSubmit}>
      <input value={form.title} onChange={e => setForm({...form, title: e.target.value})} placeholder="Title" required />
      <input value={form.brand} onChange={e => setForm({...form, brand: e.target.value})} placeholder="Brand" required />
      <input value={form.model} onChange={e => setForm({...form, model: e.target.value})} placeholder="Model" required />
      <input type="number" value={form.price} onChange={e => setForm({...form, price: e.target.value})} placeholder="Price" required />
      <input type="number" value={form.year} onChange={e => setForm({...form, year: e.target.value})} placeholder="Year" required />
      <input value={form.city} onChange={e => setForm({...form, city: e.target.value})} placeholder="City" required />
      <input value={form.contactName} onChange={e => setForm({...form, contactName: e.target.value})} placeholder="Contact Name" required />
      <input type="email" value={form.email} onChange={e => setForm({...form, email: e.target.value})} placeholder="Email" required />
      <input type="file" multiple accept="image/*" onChange={e => setImages(Array.from(e.target.files))} />
      <button type="submit">Create Vehicle</button>
    </form>
  );
};
```

---

## 🧪 Test with cURL

```bash
curl -X POST http://localhost:8080/api/user/vehicles \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "title=Test Vehicle" \
  -F "category=autoturisme" \
  -F "price=10000" \
  -F "currency=lei" \
  -F "person_type=persoana_fizica" \
  -F "brand=Test" \
  -F "model=Test" \
  -F "fuel_type=benzina" \
  -F "body_type=sedan" \
  -F "year=2020" \
  -F "condition=utilizat" \
  -F "transmission=automata" \
  -F "steering=stanga" \
  -F "city=Test City" \
  -F "contact_name=Test User" \
  -F "email=test@test.com"
```

---

## 📚 Full Documentation

See `FRONTEND_INTEGRATION_GUIDE.md` for complete documentation with:
- Detailed examples for React, Vue, Vanilla JS
- Image upload implementation
- Advanced error handling
- Validation examples
- Postman collection

---

**Endpoint**: `POST /api/user/vehicles`
**Auth**: Required (Bearer Token)
**Content-Type**: `multipart/form-data`
**Swagger**: http://localhost:8080/swagger/index.html
