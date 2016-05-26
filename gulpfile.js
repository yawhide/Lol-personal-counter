var gulp = require('gulp');
var concat = require('gulp-concat');
var rename = require('gulp-rename');
var uglify = require('gulp-uglify');
var cleanCSS = require('gulp-clean-css');

gulp.task('default', function () {
   // Your default task
});

gulp.task('scripts', function() {
    return gulp.src(['./static/js/promise.min.js', './static/js/awesomplete.min.js', './static/js/fetch.min.js', './static/js/index.js'])
        .pipe(concat('scripts.js'))
        .pipe(rename('scripts.min.js'))
        .pipe(uglify())
        .pipe(gulp.dest('static/dist/js'))
})

gulp.task('css', function(){
  return gulp.src(['./static/css/awesomplete.css', './static/css/main.css'])
      .pipe(concat('style.css'))
      .pipe(rename('style.min.css'))
      .pipe(cleanCSS())
      .pipe(gulp.dest('static/dist/css'))
})

gulp.task('concat-css', function() {
  return gulp.src(['./static/css/bootstrap.min.css', './static/dist/css/style.min.css'])
      .pipe(concat('style.min.css'))
      .pipe(gulp.dest('static/dist/css'))
})
