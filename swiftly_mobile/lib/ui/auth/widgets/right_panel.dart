import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';

class RightPanel extends StatelessWidget {
  const RightPanel({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.only(left: 16),
      child: SvgPicture.asset('assets/test.svg'),
    );
  }
}
