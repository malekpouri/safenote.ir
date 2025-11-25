# SafeNote | سیف‌نوت

SafeNote is a secure, open-source web application for sharing encrypted, self-destructing notes. It ensures your sensitive information remains private by encrypting data on the client side before it reaches the server.

سیف‌نوت یک برنامه وب امن و متن‌باز برای اشتراک‌گذاری یادداشت‌های رمزگذاری شده و خود-تخریب‌گر است. این برنامه با رمزگذاری داده‌ها در سمت کاربر (Client-side) قبل از ارسال به سرور، حریم خصوصی اطلاعات حساس شما را تضمین می‌کند.

## Features | ویژگی‌ها

### English

- **Client-Side Encryption:** Notes are encrypted in your browser using WebAssembly (WASM) before being sent to the server. The server never sees the plain text.
- **Self-Destructing Notes:** Set notes to automatically delete after being read or after a specific time duration.
- **Password Protection:** Add an extra layer of security with a password.
- **No Registration:** Use the service instantly without creating an account.
- **Internationalization:** Full support for English and Persian (Farsi) languages with RTL layout.
- **Responsive Design:** Optimized for desktop and mobile devices.
- **Character Limit:** Input limit of 10,000 characters with visual feedback.
- **Automatic Cleanup:** Expired notes are automatically deleted from the database via a background scheduler.
- **Admin Dashboard:** A public statistics page to view site activity (active and total notes).

### فارسی

- **رمزگذاری سمت کاربر:** یادداشت‌ها قبل از ارسال به سرور، در مرورگر شما با استفاده از WebAssembly رمزگذاری می‌شوند. سرور هرگز متن اصلی را نمی‌بیند.
- **یادداشت‌های خود-تخریب‌گر:** تنظیم یادداشت‌ها برای حذف خودکار پس از خوانده شدن یا پس از یک مدت زمان مشخص.
- **محافظت با رمز عبور:** افزودن یک لایه امنیتی اضافی با رمز عبور.
- **بدون نیاز به ثبت نام:** استفاده فوری از سرویس بدون نیاز به ساخت حساب کاربری.
- **چندزبانه:** پشتیبانی کامل از زبان‌های انگلیسی و فارسی با چیدمان راست‌چین (RTL).
- **طراحی واکنش‌گرا:** بهینه‌شده برای دستگاه‌های دسکتاپ و موبایل.
- **محدودیت کاراکتر:** محدودیت ورودی ۱۰,۰۰۰ کاراکتر با نمایشگر تعداد.
- **پاکسازی خودکار:** یادداشت‌های منقضی شده به صورت خودکار توسط زمان‌بند پس‌زمینه از پایگاه داده حذف می‌شوند.
- **داشبورد مدیریت:** صفحه آمار عمومی برای مشاهده فعالیت سایت (یادداشت‌های فعال و کل یادداشت‌ها).

## Tech Stack | تکنولوژی‌های استفاده شده

- **Frontend:** SvelteKit, TailwindCSS, TypeScript
- **Backend:** Go (Golang), SQLite
- **Encryption:** Go (compiled to WebAssembly for client-side execution)
- **Containerization:** Docker, Docker Compose

## Getting Started | راهنمای شروع

### Prerequisites | پیش‌نیازها

- Docker
- Docker Compose

### Installation & Running | نصب و اجرا

1. **Clone the repository | مخزن را کلون کنید:**

   ```bash
   git clone https://github.com/malekpouri/safenote.ir.git
   cd safenote.ir
   ```

2. **Run with Docker Compose | اجرا با داکر کامپوز:**

   ```bash
   docker-compose up --build
   ```

3. **Access the application | دسترسی به برنامه:**
   Open your browser and navigate to `http://localhost:3000`.

   مرورگر خود را باز کنید و به آدرس `http://localhost:3000` بروید.

## Development | توسعه

### Backend

The backend runs on port 8080.

```bash
cd backend
go run cmd/server/main.go
```

### Frontend

The frontend runs on port 3000 (or 5173 for dev).

```bash
cd frontend
npm install
npm run dev
```

## License | مجوز

This project is licensed under the MIT License.
این پروژه تحت مجوز MIT منتشر شده است.
