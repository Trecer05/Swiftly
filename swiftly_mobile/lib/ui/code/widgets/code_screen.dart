import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart' show kIsWeb;

class CodeScreen extends StatelessWidget {
  const CodeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    if (kIsWeb) {
      return const Scaffold(
        body: Center(child: Text('Download desktop version!')),
      );
    } else {
      return const Scaffold(body: Center(child: Text('Code Screen')));
    }
  }
}
