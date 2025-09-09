import 'package:flutter/material.dart';

import 'dart:io' show Platform;

class ResponsiveLayout extends StatelessWidget {
  final Widget mobile;
  final Widget desktop;
  const ResponsiveLayout({
    super.key,
    required this.mobile,
    required this.desktop,
  });

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (_, constraints) {
        if (Platform.isIOS || Platform.isAndroid) {
          return mobile;
        } else {
          return desktop;
        }
      },
    );
  }
}
