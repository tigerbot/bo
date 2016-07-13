(function () {
	'use strict';
	var gulp    = require('gulp');
	var pug     = require('gulp-pug');
	var less    = require('gulp-less');
	var browser = require('gulp-browser').browserify;

	gulp.task('html', function build_html() {
		return gulp.src('pug/*.pug')
			.pipe(pug({}))
			.pipe(gulp.dest('../public'))
		;
	});

	gulp.task('style', function build_html() {
		return gulp.src('less/*.less')
			.pipe(less({}))
			.pipe(gulp.dest('../public'))
		;
	});

	gulp.task('javascript', function combine_js() {
		return gulp.src('js/index.js')
			.pipe(browser())
			.pipe(gulp.dest('../public'))
		;
	});

	gulp.task('static', function copy_static() {
		return gulp.src('static/**')
			.pipe(gulp.dest('../public'))
		;
	});

	gulp.task('default', function () {
		gulp.start('html', 'style', 'javascript', 'static');
	});
	gulp.task('watch', function () {
		gulp.watch('pug/*.pug', ['html']);
		gulp.watch('less/*.less', ['style']);
		gulp.watch('js/*.js', ['javascript']);
		gulp.watch('static/**', ['static']);
	});
})();
