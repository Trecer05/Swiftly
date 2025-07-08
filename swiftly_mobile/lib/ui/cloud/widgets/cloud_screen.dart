import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/cloud/widgets/left_panel.dart';

class CloudScreen extends StatelessWidget {
  const CloudScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color.fromARGB(255, 9, 30, 114),
      body: Row(
        children: [
          const SizedBox(
            width: 200,
            child: LeftPanel(),
          ),
          Expanded(
            child: Container(
              color: Colors.white,
              child: const Center(child: Text('Основное содержимое')),
            ),
          ),
        ],
      ),
    );
  }
}