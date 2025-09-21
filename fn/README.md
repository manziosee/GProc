# GProc Enterprise Dashboard

> Modern, professional web interface for GProc - Enterprise Process Manager

![GProc Dashboard](https://img.shields.io/badge/Vue.js-3.0-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-5.0-3178C6?style=for-the-badge&logo=typescript&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-5.0-646CFF?style=for-the-badge&logo=vite&logoColor=white)

## ✨ Features

- **🎨 Modern UI**: Clean, professional interface inspired by PM2 and Nomad
- **🌙 Dark Mode**: Accessible dark/light theme with high contrast
- **📱 Responsive**: Mobile-first design that works on all devices
- **⚡ Real-time**: Live process monitoring and log streaming
- **🔧 Process Management**: Start, stop, restart, and monitor processes
- **📊 Analytics**: Resource usage charts and performance metrics
- **👥 User Management**: Role-based access control and user administration
- **⏰ Scheduler**: Cron job management with visual interface
- **📋 CLI Reference**: Complete command documentation
- **⚙️ Configuration**: Visual config editor with YAML/JSON support

## 🚀 Quick Start

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## 🛠️ Tech Stack

- **Framework**: Vue.js 3 with Composition API
- **Language**: TypeScript for type safety
- **Build Tool**: Vite for fast development
- **UI Library**: Naive UI for components
- **Icons**: Lucide Vue for modern icons
- **Charts**: Chart.js for data visualization
- **Styling**: CSS3 with custom properties

## 📁 Project Structure

```
fn/
├── public/
│   ├── logo.svg              # GProc brand logo
│   └── index.html           # HTML template
├── src/
│   ├── components/          # Reusable components
│   │   ├── layout/         # Layout components
│   │   ├── process/        # Process management
│   │   ├── users/          # User management
│   │   └── scheduler/      # Task scheduling
│   ├── pages/              # Main application pages
│   │   ├── Dashboard.vue   # Main dashboard
│   │   ├── ProcessManagement.vue
│   │   ├── UserManagement.vue
│   │   ├── Settings.vue
│   │   ├── CLI.vue         # CLI documentation
│   │   └── ...
│   ├── App.vue             # Root component
│   ├── main.ts             # Application entry
│   └── style.css           # Global styles
├── package.json            # Dependencies
├── vite.config.ts          # Vite configuration
└── README.md               # This file
```

## 🎨 Design System

### Colors
- **Primary**: `#10b981` (Emerald Green)
- **Secondary**: `#3b82f6` (Blue)
- **Accent**: `#8b5cf6` (Purple)
- **Success**: `#10b981` (Green)
- **Warning**: `#f59e0b` (Amber)
- **Error**: `#ef4444` (Red)

### Typography
- **Font Family**: System fonts (-apple-system, BlinkMacSystemFont, 'Segoe UI')
- **Headings**: 700 weight, tight letter spacing
- **Body**: 400 weight, 1.5 line height
- **Code**: Consolas, Monaco, monospace

## 📱 Pages Overview

### Dashboard
- Real-time process metrics
- Activity feed
- System health overview
- Quick actions

### Process Management
- Process list with status
- Start/stop/restart controls
- Resource monitoring
- Health checks

### User Management
- User profiles with avatars
- Role-based permissions
- Activity tracking
- Bulk operations

### Configuration
- YAML/JSON editor
- Environment variables
- Configuration templates
- Import/export

### Scheduler
- Cron job management
- Visual schedule builder
- Execution history
- Task templates

### CLI Reference
- Complete command documentation
- Copy-to-clipboard examples
- Categorized commands
- Quick reference guide

## 🌙 Theme System

The dashboard supports both light and dark themes with:
- **High contrast text** for accessibility
- **Consistent color variables** across components
- **Smooth transitions** between themes
- **System preference detection**

```css
/* Light theme */
--n-text-color: #1e293b;
--n-text-color-2: #64748b;

/* Dark theme */
--n-text-color: #ffffff;
--n-text-color-2: #e2e8f0;
```

## 🔧 Development

### Prerequisites
- Node.js 18+ 
- npm or yarn

### Environment Setup
```bash
# Clone the repository
git clone https://github.com/manziosee/GProc.git
cd GProc/fn

# Install dependencies
npm install

# Start development server
npm run dev
```

### Available Scripts
- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint (if configured)

## 🚀 Production Deployment

```bash
# Build the application
npm run build

# The dist/ folder contains the built application
# Deploy to your web server or CDN
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is part of GProc - Enterprise Process Manager.

## 🔗 Related

- [GProc Main Repository](https://github.com/manziosee/GProc)
- [GProc CLI Documentation](../README.md)

---

**Built with ❤️ by [Manzi Osee](mailto:manziosee3@gmail.com)**