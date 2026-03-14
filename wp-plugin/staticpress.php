<?php
/**
 * Plugin Name: StaticPress
 * Description: Provides API key authentication for StaticPress CLI
 * Version: 1.0.0
 * Author: StaticPress
 * License: MIT
 */

if (!defined('ABSPATH')) {
    exit;
}

define('STATICPRESS_VERSION', '1.0.0');
define('STATICPRESS_PLUGIN_DIR', plugin_dir_path(__FILE__));

require_once STATICPRESS_PLUGIN_DIR . 'includes/class-staticpress-api.php';
require_once STATICPRESS_PLUGIN_DIR . 'includes/class-staticpress-settings.php';

class StaticPress {
    private static $instance = null;

    public static function get_instance() {
        if (null === self::$instance) {
            self::$instance = new self();
        }
        return self::$instance;
    }

    private function __construct() {
        add_action('rest_api_init', array($this, 'register_routes'));
        add_action('admin_menu', array($this, 'add_admin_menu'));
        add_action('admin_init', array($this, 'register_settings'));
    }

    public function register_routes() {
        StaticPress_API::register_routes();
    }

    public function add_admin_menu() {
        add_options_page(
            'StaticPress',
            'StaticPress',
            'manage_options',
            'staticpress',
            array('StaticPress_Settings', 'render_page')
        );
    }

    public function register_settings() {
        StaticPress_Settings::register_settings();
    }
}

function staticpress() {
    return StaticPress::get_instance();
}

staticpress();
